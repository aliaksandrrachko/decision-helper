package query

type Expression interface {
	Renderable
	IsNull() Expression
	IsNotNull() Expression
	In(values []any) Expression
	Eq(propertyName string, value any) Expression
}

type Renderable interface {
	Render() string
	//Render(RenderingContext renderingContext) string
}
