// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package render

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/mr55p-dev/pagemail/db/queries"
	"github.com/mr55p-dev/pagemail/render/icons"
	"github.com/mr55p-dev/pagemail/render/wrapper"
)

func App(user *queries.User, pages []queries.Page) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"app-header\"><form hx-post=\"/app/page\" hx-target=\"#pages-list\" hx-swap=\"afterbegin\" hx-target-error=\"#messages\"><fieldset role=\"group\"><input placeholder=\"Page URL\" autocomplete=\"off\" type=\"url\" id=\"page-input\" name=\"url\"><div class=\"button-group\"><button type=\"button\" onclick=\"pasteContents(&#39;page-input&#39;)\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = icons.Clipboard().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button> <button type=\"submit\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = icons.Save().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div></fieldset></form></div><form class=\"tw-group/pages\" id=\"pages-form\"><div class=\"tw-grid-cols-2 tw-gap-2 tw-px-2 tw-pt-4 tw-grid group-data-[selecting=true]/pages:tw-hidden\"><button type=\"button\" class=\"tw-px-2 tw-py-2 tw-text-left tw-w-full tw-text-gray-700 tw-font-semibold tw-bg-primary tw-rounded-lg tw-border tw-border-transparent hover:tw-border-brand-800 tw-transition-colors tw-ease-in-out\" onclick=\"\">Search</button> <button type=\"button\" class=\"tw-px-2 tw-py-2 tw-text-left tw-w-full tw-text-gray-700 tw-font-semibold tw-bg-primary tw-rounded-lg tw-border tw-border-transparent hover:tw-border-brand-800 tw-transition-colors tw-ease-in-out\" onclick=\"handleSelecting(this)\">Select</button></div><div class=\"tw-grid-cols-2 tw-gap-2 tw-px-2 tw-pt-4 tw-hidden group-data-[selecting=true]/pages:tw-grid\"><button type=\"button\" class=\"tw-px-2 tw-py-2 tw-text-left tw-w-full tw-text-gray-700 tw-font-semibold tw-bg-primary tw-rounded-lg tw-border tw-border-transparent hover:tw-border-brand-800 tw-transition-colors tw-ease-in-out\" onclick=\"\">Delete</button> <button type=\"button\" class=\"tw-px-2 tw-py-2 tw-text-left tw-w-full tw-text-gray-700 tw-font-semibold tw-bg-primary tw-rounded-lg tw-border tw-border-transparent hover:tw-border-brand-800 tw-transition-colors tw-ease-in-out\" onclick=\"handleSelecting(this)\">Cancel</button></div><div id=\"pages-list\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, page := range pages {
				templ_7745c5c3_Err = PageCard(page).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></form><script>\n\t\t\tfunction isSelecting() {\n\t\t\t\treturn document.querySelector(\"form#pages-form\").dataset.selecting === \"true\"\n\t\t\t}\n\n\t\t\tfunction handleSelecting(trigger) {\n\t\t\t\tconst form = trigger.closest(\"form\");\n\t\t\t\tform.dataset.selecting = form.dataset.selecting === \"true\" ? \"false\" : \"true\"\n\t\t\t\tconst selecting = form.dataset.selecting === \"true\";\n\t\t\t\tif (!selecting) {\n\t\t\t\t\tform.querySelectorAll(\"input[type=checkbox]\").forEach(box => {\n\t\t\t\t\t\tbox.checked = false\n\t\t\t\t\t});\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tfunction handleSelect(element) {\n\t\t\t\tif (!isSelecting()) {\n\t\t\t\t\treturn\n\t\t\t\t}\n\t\t\t\tconst id = element.dataset.id;\n\t\t\t\tconst input = element.querySelector(\"input[type=checkbox]\");\n\t\t\t\tinput.checked = !input.checked;\n\t\t\t}\n\t\t</script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = wrapper.New(user, wrapper.WithTitle("Pages")).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
