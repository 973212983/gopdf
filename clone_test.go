package gopdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestClonePreservesCharacterToGlyphIndex(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	// Build a PDF with text
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.Cell(nil, "Hello World 123 测试中文")

	// Inspect original's CharacterToGlyphIndex before clone
	t.Log("=== Original objects before clone ===")
	for i, obj := range pdf.pdfObjs {
		switch o := obj.(type) {
		case *CIDFontObj:
			if o.PtrToSubsetFontObj != nil {
				vals := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				keys := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
				t.Logf("  Original CIDFontObj[%d]: %d glyphs, keys=%v", i, len(vals), keys)
			}
		case *SubsetFontObj:
			if o.CharacterToGlyphIndex != nil {
				vals := o.CharacterToGlyphIndex.AllVals()
				t.Logf("  Original SubsetFontObj[%d]: %d glyphs", i, len(vals))
			}
		}
	}

	// Clone
	clone := pdf.Clone()

	// Inspect clone's CharacterToGlyphIndex after clone
	t.Log("=== Clone objects after clone ===")
	for i, obj := range clone.pdfObjs {
		switch o := obj.(type) {
		case *CIDFontObj:
			if o.PtrToSubsetFontObj != nil {
				vals := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				keys := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
				t.Logf("  Clone CIDFontObj[%d]: %d glyphs, keys=%v", i, len(vals), keys)
				if len(vals) == 0 {
					t.Errorf("Clone CIDFontObj[%d] has EMPTY CharacterToGlyphIndex!", i)
				}
			} else {
				t.Errorf("Clone CIDFontObj[%d] has nil PtrToSubsetFontObj", i)
			}
		case *SubsetFontObj:
			if o.CharacterToGlyphIndex != nil {
				vals := o.CharacterToGlyphIndex.AllVals()
				t.Logf("  Clone SubsetFontObj[%d]: %d glyphs", i, len(vals))
				if len(vals) == 0 {
					t.Errorf("Clone SubsetFontObj[%d] has EMPTY CharacterToGlyphIndex!", i)
				}
			}
		case *UnicodeMap:
			if o.PtrToSubsetFontObj != nil {
				keys := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
				t.Logf("  Clone UnicodeMap[%d]: %d keys", i, len(keys))
			}
		}
	}

	// Write both to buffers and compare
	var origBuf, cloneBuf bytes.Buffer
	_, err = pdf.WriteTo(&origBuf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = clone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}

	// Write to files for manual comparison
	os.WriteFile("./test/out/original_clone_test.pdf", origBuf.Bytes(), 0644)
	os.WriteFile("./test/out/clone_clone_test.pdf", cloneBuf.Bytes(), 0644)

	// Check for W arrays in both outputs
	origStr := origBuf.String()
	cloneStr := cloneBuf.String()

	// Count W array entries
	origWCount := strings.Count(origStr, "/W [")
	cloneWCount := strings.Count(cloneStr, "/W [")
	t.Logf("Original W arrays: %d, Clone W arrays: %d", origWCount, cloneWCount)

	// Compare W array content
	origWContent := extractWArrays(origStr)
	cloneWContent := extractWArrays(cloneStr)
	t.Logf("Original W: %s", origWContent)
	t.Logf("Clone W: %s", cloneWContent)

	if origWContent != cloneWContent {
		t.Errorf("W array mismatch!\n  Original: %s\n  Clone:    %s", origWContent, cloneWContent)
	}
}

func extractWArrays(pdfContent string) string {
	var result strings.Builder
	for _, line := range strings.Split(pdfContent, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "/W") {
			result.WriteString(trimmed)
			result.WriteString(" | ")
		}
	}
	return result.String()
}

func TestCloneProducesValidPDF(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	// Build a more realistic PDF with multiple operations
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}

	pdf.Cell(nil, "Hello World")
	pdf.Br(20)
	pdf.Cell(nil, "Line 2 with more text")
	pdf.Br(20)
	pdf.SetFont("LiberationSerif-Regular", "", 10)
	pdf.Cell(nil, "Smaller text line")
	pdf.Br(20)
	pdf.SetFont("LiberationSerif-Regular", "", 16)
	pdf.Cell(nil, "Larger text line")

	// Clone and write
	clone := pdf.Clone()

	// Also add a new page to clone and write more text
	clone.AddPage()
	clone.Cell(nil, "Added after clone")

	var origBuf, cloneBuf bytes.Buffer
	_, err = pdf.WriteTo(&origBuf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = clone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("./test/out/original_multi.pdf", origBuf.Bytes(), 0644)
	os.WriteFile("./test/out/clone_multi.pdf", cloneBuf.Bytes(), 0644)

	t.Logf("Original PDF size: %d bytes", origBuf.Len())
	t.Logf("Clone PDF size: %d bytes", cloneBuf.Len())

	// Verify both are valid PDFs
	for name, data := range map[string][]byte{
		"original": origBuf.Bytes(),
		"clone":    cloneBuf.Bytes(),
	} {
		if !bytes.HasPrefix(data, []byte("%PDF-1.7")) {
			t.Errorf("%s: invalid PDF header", name)
		}
		if !bytes.Contains(data, []byte("%%EOF")) {
			t.Errorf("%s: missing %%EOF", name)
		}

		str := string(data)
		wCount := strings.Count(str, "/W [")
		if wCount == 0 {
			t.Errorf("%s: no W arrays found!", name)
		}
		t.Logf("%s: %d W arrays, size=%d", name, wCount, len(data))
	}

	// Compare W array content
	origW := extractWArrays(string(origBuf.Bytes()))
	cloneW := extractWArrays(string(cloneBuf.Bytes()))
	t.Logf("Original W arrays: %s", origW)
	t.Logf("Clone W arrays: %s", cloneW)
}

func TestGoPdfClone(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	// Test from the user's perspective - simple case
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	pdf.Cell(nil, "Hello World")

	// Clone immediately after text
	pdfClone := pdf.Clone()

	// Check clone internals
	for _, obj := range pdfClone.pdfObjs {
		if cidFont, ok := obj.(*CIDFontObj); ok {
			if cidFont.PtrToSubsetFontObj != nil {
				vals := cidFont.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				if len(vals) == 0 {
					t.Error("BUG: Clone CIDFontObj has empty CharacterToGlyphIndex.AllVals()")
				} else {
					t.Logf("OK: Clone CIDFontObj has %d glyph indices", len(vals))
				}
			}
		}
	}

	// WritePdf on clone
	var cloneBuf bytes.Buffer
	_, err = pdfClone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}
	cloneStr := string(cloneBuf.Bytes())

	// Check for populated W arrays in output
	for _, line := range strings.Split(cloneStr, "\n") {
		if strings.Contains(line, "/W [") {
			if strings.Contains(line, "/W []") {
				t.Errorf("BUG: Clone PDF has EMPTY W array in output!")
			} else {
				t.Logf("OK: Clone PDF has populated W array: %.100s", strings.TrimSpace(line))
			}
		}
	}
}

func TestCloneAfterMultiplePages(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}

	// Multiple pages with text
	for page := 0; page < 3; page++ {
		if page > 0 {
			pdf.AddPage()
		}
		for i := 0; i < 5; i++ {
			pdf.Cell(nil, fmt.Sprintf("Page %d Line %d: Hello World 测试中文", page+1, i+1))
			pdf.Br(20)
		}
	}

	// Clone
	clone := pdf.Clone()

	// Verify clone's CIDFont data
	for i, obj := range clone.pdfObjs {
		if cidFont, ok := obj.(*CIDFontObj); ok {
			if cidFont.PtrToSubsetFontObj != nil {
				vals := cidFont.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				t.Logf("Clone CIDFontObj[%d]: %d glyphs", i, len(vals))
				if len(vals) == 0 {
					t.Errorf("BUG: Clone CIDFontObj[%d] empty!", i)
				}
			}
		}
	}

	// WritePdf clone
	var cloneBuf bytes.Buffer
	_, err = clone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}

	cloneStr := string(cloneBuf.Bytes())
	for _, line := range strings.Split(cloneStr, "\n") {
		if strings.Contains(line, "/W [") && strings.Contains(line, "/W []") {
			// Find which object this is
			t.Errorf("BUG: Clone output contains empty W array!")
		}
	}
}

func TestCloneAddNewCharsAfterClone(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}

	// Add text with characters A, B, C
	pdf.Cell(nil, "ABC")

	// Clone - now all 5 font objects + curr have independent SubsetFontObj copies
	clone := pdf.Clone()

	// Switch to a different font in clone (to ensure curr.FontISubset changes)
	err = clone.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}

	// Add text with NEW characters D, E, F (not in original)
	// After this, clone.curr.FontISubset has A,B,C,D,E,F
	// But clone's CIDFontObj.PtrToSubsetFontObj only has A,B,C
	clone.Cell(nil, "DEF")

	// Now check: does the clone's CIDFontObj have all 6 characters?
	t.Log("=== Checking clone's CIDFontObj after adding new chars ===")
	for i, obj := range clone.pdfObjs {
		if cidFont, ok := obj.(*CIDFontObj); ok {
			if cidFont.PtrToSubsetFontObj != nil {
				vals := cidFont.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				keys := cidFont.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
				t.Logf("  Clone CIDFontObj[%d]: %d glyphs, keys=%v", i, len(vals), keys)
				// Should have at least ABC + DEF = 6 chars (+ space = 7)
				if len(vals) < 6 {
					t.Errorf("BUG: Clone CIDFontObj only has %d glyphs, expected 6+ (characters from both before and after clone)", len(vals))
				}
			}
		}
	}

	// Write both and compare W arrays
	var origBuf, cloneBuf bytes.Buffer
	_, err = pdf.WriteTo(&origBuf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = clone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}

	origStr := origBuf.String()
	cloneStr := cloneBuf.String()

	origW := extractWArrays(origStr)
	cloneW := extractWArrays(cloneStr)
	t.Logf("Original W: %s", origW)
	t.Logf("Clone W: %s", cloneW)

	// Count entries in W arrays
	origCount := strings.Count(origW, "[")
	cloneCount := strings.Count(cloneW, "[")
	t.Logf("Original W entries: %d, Clone W entries: %d", origCount, cloneCount)

	// The clone should have MORE or EQUAL W entries (it has ABC + DEF)
	if cloneCount < origCount {
		t.Errorf("BUG: Clone W array has FEWER entries than original (%d < %d)! New chars added after clone are missing!", cloneCount, origCount)
	}

	// Specific check: clone should have both ABC and DEF glyph entries
	// W arrays use glyph/CID indices (not character rune codes), so we
	// look up each character's glyph index from CharacterToGlyphIndex.
	for _, obj := range clone.pdfObjs {
		if cidFont, ok := obj.(*CIDFontObj); ok && cidFont.PtrToSubsetFontObj != nil {
			for _, ch := range "ABCDEF" {
				gid, err := cidFont.PtrToSubsetFontObj.CharCodeToGlyphIndex(ch)
				if err != nil {
					t.Errorf("BUG: no glyph index for '%c' (rune %d): %v", ch, ch, err)
					continue
				}
				expectedStr := fmt.Sprintf("%d[", gid)
				if !strings.Contains(cloneStr, expectedStr) {
					t.Errorf("BUG: Clone W array missing glyph %d for '%c'!", gid, ch)
				} else {
					t.Logf("OK: glyph %d ('%c') found in W array as %q", gid, ch, expectedStr)
				}
			}
		}
	}

	t.Log("OK: All expected glyphs found in clone W array")
}

func TestCloneMultipleFontsAddNewChars(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Fatal(err)
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}

	// Add text with some characters
	pdf.Cell(nil, "ABC")

	clone := pdf.Clone()

	// Add more characters to clone (different from original)
	err = clone.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Fatal(err)
	}
	clone.Cell(nil, "DEF")

	// Check clone's CIDFontObj — should have all 6 chars + space (7 total)
	t.Log("=== Multi-font clone check ===")
	for i, obj := range clone.pdfObjs {
		switch o := obj.(type) {
		case *CIDFontObj:
			if o.PtrToSubsetFontObj != nil {
				vals := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
				keys := o.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
				t.Logf("  Clone CIDFontObj[%d]: %d glyphs, keys=%v", i, len(vals), keys)
				if len(vals) < 6 {
					t.Errorf("BUG: Clone CIDFontObj[%d] only has %d glyphs (expected 6+)!", i, len(vals))
				}
			}
		}
	}

	// Write clone and check W arrays are populated
	var cloneBuf bytes.Buffer
	_, err = clone.WriteTo(&cloneBuf)
	if err != nil {
		t.Fatal(err)
	}

	cloneStr := string(cloneBuf.Bytes())
	hasEmptyW := false
	for _, line := range strings.Split(cloneStr, "\n") {
		if strings.Contains(line, "/W []") {
			hasEmptyW = true
			t.Errorf("BUG: Clone has EMPTY W array!")
		}
	}
	if !hasEmptyW {
		t.Log("OK: No empty W arrays in clone output")
	}

	// Verify all 6 characters have glyph entries in clone W array via CID mapping
	for _, obj := range clone.pdfObjs {
		if cidFont, ok := obj.(*CIDFontObj); ok && cidFont.PtrToSubsetFontObj != nil {
			for _, ch := range "ABCDEF" {
				gid, err := cidFont.PtrToSubsetFontObj.CharCodeToGlyphIndex(ch)
				if err != nil {
					t.Errorf("BUG: no glyph index for '%c' in clone: %v", ch, err)
					continue
				}
				expectedStr := fmt.Sprintf("%d[", gid)
				if !strings.Contains(cloneStr, expectedStr) {
					t.Errorf("BUG: Clone W array missing glyph %d for '%c'!", gid, ch)
				}
			}
			t.Log("OK: All pre-clone (ABC) and post-clone (DEF) characters have glyph entries")
		}
	}

	os.WriteFile("./test/out/clone_multi_font.pdf", cloneBuf.Bytes(), 0644)
	t.Logf("Clone multi-font PDF: %d bytes", cloneBuf.Len())
}
