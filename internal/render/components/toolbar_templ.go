// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Toolbar() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form class=\"\n			mx-auto\n			bg-grey-800 rounded-lg\n			p-8\n			grid gap-4 grid-columns-4\n			md:grid-columns-3 md:grid-rows-2\n		\"><input id=\"save-page-input\" placeholder=\"URL\" type=\"url\" name=\"url\" class=\"md:col-span-3\" hx-trigger=\"submit,click from:#save-page-submit\" hx-post=\"/pages/\" hx-target=\"#pages\" hx-swap=\"afterbegin\"> <button id=\"add-page-button\" class=\"border btn-slim fill-primary hocus:hollow-primary\" type=\"submit\">Add page</button> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, PasteFromClipboard("save-page-input"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"btn-slim fill-secondary hocus:hollow-secondary\" onClick=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = PasteFromClipboard("save-page-input")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Paste from clipboard</button> <button class=\"btn-slim fill-bad hocus:hollow-bad\">Delete all</button></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func PasteFromClipboard(toId string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_PasteFromClipboard_5668`,
		Function: `function __templ_PasteFromClipboard_5668(toId){const elem = document.getElementById(toId)
	console.log("Grabbed", elem)
	navigator.clipboard.readText().then(txt => elem.setAttribute("value", txt)).catch(err => console.error("Failed to paste", error))
}`,
		Call:       templ.SafeScript(`__templ_PasteFromClipboard_5668`, toId),
		CallInline: templ.SafeScriptInline(`__templ_PasteFromClipboard_5668`, toId),
	}
}
