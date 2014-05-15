package main

import (
	"fmt"
	"strings"
)

type Type struct {
	Name         string // Name of the type without modifiers
	PointerLevel int    // Number of levels of declared indirection to the type
	CDefinition  string // Raw C definition
}

type Typedef struct {
	Name        string // Name of the defined type (or included types)
	CDefinition string // Raw C definition
}

func (t Type) String() string {
	return fmt.Sprintf("%s%s [%s]", t.Name, t.pointers(), t.CDefinition)
}

func (t Type) pointers() string {
	return strings.Repeat("*", t.PointerLevel)
}

func (t Type) IsVoid() bool {
	return (t.Name == "void" || t.Name == "GLvoid") && t.PointerLevel == 0
}

// CType returns the C definition of the type.
func (t Type) CType() string {
	return t.CDefinition
}

// GoType returns the Go definition of the type.
func (t Type) GoType() string {
	switch t.Name {
	case "GLbyte":
		return t.pointers() + "int8"
	case "GLubyte":
		return t.pointers() + "uint8"
	case "GLshort":
		return t.pointers() + "int16"
	case "GLushort":
		return t.pointers() + "uint16"
	case "GLint":
		return t.pointers() + "int32"
	case "GLuint":
		return t.pointers() + "uint32"
	case "GLint64", "GLint64EXT":
		return t.pointers() + "int64"
	case "GLuint64", "GLuint64EXT":
		return t.pointers() + "uint64"
	case "GLfloat", "GLclampf":
		return t.pointers() + "float32"
	case "GLdouble", "GLclampd":
		return t.pointers() + "float64"
	case "GLclampx":
		return t.pointers() + "int32"
	case "GLsizei":
		return t.pointers() + "int32"
	case "GLfixed":
		return t.pointers() + "int32"
	case "GLchar", "GLcharARB":
		return t.pointers() + "int8"
	case "GLboolean":
		return t.pointers() + "bool"
	case "GLenum", "GLbitfield":
		return t.pointers() + "uint32"
	case "GLhalf", "GLhalfNV": // Go has no 16-bit floating point type
		return t.pointers() + "uint16"
	case "void", "GLvoid":
		if t.PointerLevel == 1 {
			return "uintptr"
		} else if t.PointerLevel == 2 {
			return "*uintptr"
		}
	case "GLintptr", "GLintptrARB":
		return t.pointers() + "int"
	case "GLsizeiptr", "GLsizeiptrARB":
		return t.pointers() + "int"
	case "GLhandleARB", "GLeglImagesOES", "GLvdpauSurfaceARB":
		return t.pointers() + "uintptr"
	case "GLsync":
		return t.pointers() + "unsafe.Pointer"
	case "GLDEBUGPROC":
		return "unsafe.Pointer"
	}
	return t.pointers() + "C." + t.Name
}

// ConvertGoToC returns an expression that converts a variable from the Go type to the C type.
func (t Type) ConvertGoToC(name string) string {
	switch t.Name {
	case "GLboolean":
		if t.PointerLevel == 0 {
			return fmt.Sprintf("(C.GLboolean)(boolToInt(%s))", name)
		}
	case "void", "GLvoid":
		if t.PointerLevel == 1 {
			return fmt.Sprintf("unsafe.Pointer(%s)", name)
		} else if t.PointerLevel == 2 {
			return fmt.Sprintf("(*unsafe.Pointer)(unsafe.Pointer(%s))", name)
		}
	}
	if t.PointerLevel >= 1 {
		return fmt.Sprintf("(%sC.%s)(unsafe.Pointer(%s))", t.pointers(), t.Name, name)
	}
	return fmt.Sprintf("(%sC.%s)(%s)", t.pointers(), t.Name, name)
}

// ConvertCToGo converts from the C type to the Go type.
func (t Type) ConvertCToGo(name string) string {
	if t.Name == "GLboolean" {
		return fmt.Sprintf("%s == TRUE", name)
	}
	return fmt.Sprintf("(%s)(%s)", t.GoType(), name)
}
