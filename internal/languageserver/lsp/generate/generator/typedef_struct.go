package generator

import "github.com/dogmatiq/dogmacli/internal/languageserver/lsp/generate/model"

func (g typeDefGen) Struct(d model.Struct) {
}

// // VisitEnumeration declares a Go struct type.
// func (g *typeDefGenerator) VisitStructure(def metamodel.Struct) {
// 	// documentation(g, t.Documentation)

// 	// g.
// 	// 	Type().
// 	// 	Id(normalized(t.Name)).
// 	// 	StructFunc(func(g *jen.Group) {
// 	// 		// 			for _, t := range t.EmbeddedTypes {
// 	// 		// 				g.Add(getTypeInfo(t).Code)
// 	// 		// 			}

// 	// 		// 			if len(t.EmbeddedTypes) > 0 && len(t.Properties) > 0 {
// 	// 		// 				g.Line()
// 	// 		// 			}

// 	// 		// 			for _, p := range t.Properties {
// 	// 		// 				documentation(g, p.Documentation)
// 	// 		// 				structureProperty(g, p)
// 	// 		// 			}
// 	// 	})

// 	// 	g.
// 	// 		Func().
// 	// 		Params(
// 	// 			jen.Id("x").Id(typeName),
// 	// 		).
// 	// 		Id("encode").
// 	// 		Params(
// 	// 			jen.Id("w").Op("*").Qual("bytes", "Buffer"),
// 	// 		).
// 	// 		Params(
// 	// 			jen.Error(),
// 	// 		).
// 	// 		BlockFunc(func(g *jen.Group) {
// 	// 			g.
// 	// 				Id("w").Dot("WriteByte").
// 	// 				Call(
// 	// 					jen.LitRune('{'),
// 	// 				)

// 	// 			for _, p := range t.Properties {
// 	// 				encodeStructureProperty(g, p)
// 	// 			}

// 	// 			g.
// 	// 				Id("w").Dot("WriteByte").
// 	// 				Call(
// 	// 					jen.LitRune('}'),
// 	// 				)

// 	// 			g.Return(jen.Nil())
// 	// 		})
// 	// }

// 	// // structureProperty declares a Go struct property.
// 	// func structureProperty(
// 	// 	g *jen.Group,
// 	// 	p metamodel.StructureProperty,
// 	// ) {
// 	// 	info := getTypeInfo(p.Type)
// 	// 	code := info.Code

// 	// 	if p.Optional && info.OptionalPointer {
// 	// 		code = jen.Op("*").Add(code)
// 	// 	}

// 	// 	g.
// 	// 		Id(exported(p.Name)).
// 	// 		Add(code)
// 	// }

// 	// // encodeStructureProperty generates code to encode a Go struct property.
// 	// func encodeStructureProperty(
// 	// 	g *jen.Group,
// 	// 	p metamodel.StructureProperty,
// 	// ) {
// 	// 	info := getTypeInfo(p.Type)

// 	//	if !p.Optional {
// 	//		if info.HasEncodeMethod {
// 	//			g.
// 	//				Id("x").Dot(exported(p.Name)).Dot("encode").
// 	//				Call(
// 	//					jen.Id("w"),
// 	//				)
// 	//		} else {
// 	//			g.
// 	//				Qual("encoding/json", "NewEncoder").
// 	//				Call(
// 	//					jen.Id("w"),
// 	//				).
// 	//				Dot("Encode").
// 	//				Call(
// 	//					jen.Id("x").Dot(exported(p.Name)),
// 	//				)
// 	//		}
// 	//	} else if info.OptionalPointer {
// 	//
// 	//		g.
// 	//			If(
// 	//				jen.Id("x").Dot(exported(p.Name)).Op("!=").Nil(),
// 	//			).
// 	//			Block()
// 	//	} else {
// 	//
// 	//		g.
// 	//			If(
// 	//				jen.Id("x").Dot(exported(p.Name)).Op("!=").Add(info.Zero),
// 	//			).
// 	//			Block()
// 	//	}
// }
