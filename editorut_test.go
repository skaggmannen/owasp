package owasp

import (
	"io"
	"testing"
)

type EditorUnderTest struct {
	*testing.T
	*Editor
}

func (me *Editor) UnderTest(t *testing.T) *EditorUnderTest {
	return &EditorUnderTest{T: t, Editor: me}
}

func (me *EditorUnderTest) shouldSetVerified(id string) {
	err := me.SetVerified(id)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustSetVerified(id string) {
	err := me.SetVerified(id)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldImportFile(filename string) {
	err := me.ImportFile(filename)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustImportFile(filename string) {
	err := me.ImportFile(filename)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldImport(r io.Reader) {
	err := me.Import(r)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustImport(r io.Reader) {
	err := me.Import(r)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldTidyExport(w io.Writer) {
	err := me.TidyExport(w)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustTidyExport(w io.Writer) {
	err := me.TidyExport(w)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldExport(w io.Writer) {
	err := me.Export(w)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustExport(w io.Writer) {
	err := me.Export(w)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}

func (me *EditorUnderTest) shouldWriteReport(w io.Writer) {
	err := me.WriteReport(w)
	if err != nil {
		me.T.Helper()
		me.T.Error(err)
	}
}

func (me *EditorUnderTest) mustWriteReport(w io.Writer) {
	err := me.WriteReport(w)
	if err != nil {
		me.T.Helper()
		me.T.Fatal(err)
	}
}
