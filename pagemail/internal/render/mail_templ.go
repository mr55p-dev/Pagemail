// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package render

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/mr55p-dev/pagemail/internal/db"
	"time"
)

func MailDigest(date *time.Time, name string, pages []db.Page) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"><title>pagemail digest from ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(date.Format("02-01-2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 13, Col: 58}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><style>\n      /* -------------------------------------\n              GLOBAL RESETS\n          ------------------------------------- */\n\n      /*All the styling goes here*/\n\n      img {\n      \tborder: none;\n      \t-ms-interpolation-mode: bicubic;\n      \tmax-width: 100%;\n      }\n\n      body {\n      \tbackground-color: #f6f6f6;\n      \tfont-family: sans-serif;\n      \t-webkit-font-smoothing: antialiased;\n      \tfont-size: 14px;\n      \tline-height: 1.4;\n      \tmargin: 0;\n      \tpadding: 0;\n      \t-ms-text-size-adjust: 100%;\n      \t-webkit-text-size-adjust: 100%;\n      }\n\n      table {\n      \tborder-collapse: separate;\n      \tmso-table-lspace: 0pt;\n      \tmso-table-rspace: 0pt;\n      \twidth: 100%;\n      }\n\n      table td {\n      \tfont-family: sans-serif;\n      \tfont-size: 14px;\n      \tvertical-align: top;\n      }\n\n      /* -------------------------------------\n              BODY & CONTAINER\n          ------------------------------------- */\n\n      .body {\n      \tbackground-color: #f6f6f6;\n      \twidth: 100%;\n      }\n\n      /* Set a max-width, and make it display as block so it will automatically stretch to that width, but will also shrink down on a phone or something */\n      .container {\n      \tdisplay: block;\n      \tmargin: 0 auto !important;\n      \t/* makes it centered */\n      \tmax-width: 580px;\n      \tpadding: 10px;\n      \twidth: 580px;\n      \toverflow-x: hidden;\n      }\n\n      /* This should also be a block element, so that it will fill 100% of the .container */\n      .content {\n      \tbox-sizing: border-box;\n      \tdisplay: block;\n      \tmargin: 0 auto;\n      \tmax-width: 580px;\n      \tpadding: 10px;\n      }\n\n      /* -------------------------------------\n              HEADER, FOOTER, MAIN\n          ------------------------------------- */\n      .main {\n      \tbackground: #ffffff;\n      \tborder-radius: 3px;\n      \twidth: 100%;\n      }\n\n      .wrapper {\n      \tbox-sizing: border-box;\n      \tpadding: 20px;\n      }\n\n      .content-block {\n      \tpadding-bottom: 10px;\n      \tpadding-top: 10px;\n      }\n\n      .footer {\n      \tclear: both;\n      \tmargin-top: 10px;\n      \ttext-align: center;\n      \twidth: 100%;\n      }\n\n      .footer td,\n      .footer p,\n      .footer span,\n      .footer a {\n      \tcolor: #999999;\n      \tfont-size: 12px;\n      \ttext-align: center;\n      }\n\n      /* -------------------------------------\n              TYPOGRAPHY\n          ------------------------------------- */\n      h1,\n      h2,\n      h3,\n      h4 {\n      \tcolor: #000000;\n      \tfont-family: sans-serif;\n      \tfont-weight: 400;\n      \tline-height: 1.4;\n      \tmargin: 0;\n      \tmargin-bottom: 30px;\n      }\n\n      h1 {\n      \tfont-size: 35px;\n      \tfont-weight: 300;\n      \ttext-align: center;\n      \ttext-transform: capitalize;\n      }\n\n      p,\n      ul,\n      ol {\n      \tfont-family: sans-serif;\n      \tfont-size: 14px;\n      \tfont-weight: normal;\n      \tmargin: 0;\n      \tmargin-bottom: 15px;\n      }\n\n      p li,\n      ul li,\n      ol li {\n      \tlist-style-position: inside;\n      \tmargin-left: 5px;\n      }\n\n      a {\n      \tcolor: #3498db;\n      \ttext-decoration: underline;\n      }\n\n      /* -------------------------------------\n              BUTTONS\n          ------------------------------------- */\n      .btn {\n      \tbox-sizing: border-box;\n      \twidth: 100%;\n      }\n\n      .btn>tbody>tr>td {\n      \tpadding-bottom: 15px;\n      }\n\n      .btn table {\n      \twidth: auto;\n      }\n\n      .btn table td {\n      \tbackground-color: #ffffff;\n      \tborder-radius: 5px;\n      \ttext-align: center;\n      }\n\n      .btn a {\n      \tbackground-color: #ffffff;\n      \tborder: solid 1px #3498db;\n      \tborder-radius: 5px;\n      \tbox-sizing: border-box;\n      \tcolor: #3498db;\n      \tcursor: pointer;\n      \tdisplay: inline-block;\n      \tfont-size: 14px;\n      \tfont-weight: bold;\n      \tmargin: 0;\n      \tpadding: 12px 25px;\n      \ttext-decoration: none;\n      \ttext-transform: capitalize;\n      }\n\n      .btn-primary table td {\n      \tbackground-color: #3498db;\n      }\n\n      .btn-primary a {\n      \tbackground-color: #3498db;\n      \tborder-color: #3498db;\n      \tcolor: #ffffff;\n      }\n\n      /* -------------------------------------\n              OTHER STYLES THAT MIGHT BE USEFUL\n          ------------------------------------- */\n      .last {\n      \tmargin-bottom: 0;\n      }\n\n      .first {\n      \tmargin-top: 0;\n      }\n\n      .align-center {\n      \ttext-align: center;\n      }\n\n      .align-right {\n      \ttext-align: right;\n      }\n\n      .align-left {\n      \ttext-align: left;\n      }\n\n      .clear {\n      \tclear: both;\n      }\n\n      .mt0 {\n      \tmargin-top: 0;\n      }\n\n      .mb0 {\n      \tmargin-bottom: 0;\n      }\n\n      .preheader {\n      \tcolor: transparent;\n      \tdisplay: none;\n      \theight: 0;\n      \tmax-height: 0;\n      \tmax-width: 0;\n      \topacity: 0;\n      \toverflow: hidden;\n      \tmso-hide: all;\n      \tvisibility: hidden;\n      \twidth: 0;\n      }\n\n      .powered-by a {\n      \ttext-decoration: none;\n      }\n\n      hr {\n      \tborder: 0;\n      \tborder-bottom: 1px solid #f6f6f6;\n      \tmargin: 20px 0;\n      }\n\n      /* -------------------------------------\n              RESPONSIVE AND MOBILE FRIENDLY STYLES\n          ------------------------------------- */\n      @media only screen and (max-width: 620px) {\n      \ttable.body h1 {\n      \t\tfont-size: 28px !important;\n      \t\tmargin-bottom: 10px !important;\n      \t}\n\n      \ttable.body p,\n      \ttable.body ul,\n      \ttable.body ol,\n      \ttable.body td,\n      \ttable.body span,\n      \ttable.body a {\n      \t\tfont-size: 16px !important;\n      \t}\n\n      \ttable.body .wrapper,\n      \ttable.body .article {\n      \t\tpadding: 10px !important;\n      \t}\n\n      \ttable.body .content {\n      \t\tpadding: 0 !important;\n      \t}\n\n      \ttable.body .container {\n      \t\tpadding: 0 !important;\n      \t\twidth: 100% !important;\n      \t}\n\n      \ttable.body .main {\n      \t\tborder-left-width: 0 !important;\n      \t\tborder-radius: 0 !important;\n      \t\tborder-right-width: 0 !important;\n      \t}\n\n      \ttable.body .btn table {\n      \t\twidth: 100% !important;\n      \t}\n\n      \ttable.body .btn a {\n      \t\twidth: 100% !important;\n      \t}\n\n      \ttable.body .img-responsive {\n      \t\theight: auto !important;\n      \t\tmax-width: 100% !important;\n      \t\twidth: auto !important;\n      \t}\n      }\n\n      /* -------------------------------------\n              PRESERVE THESE STYLES IN THE HEAD\n          ------------------------------------- */\n      @media all {\n      \t.ExternalClass {\n      \t\twidth: 100%;\n      \t}\n\n      \t.ExternalClass,\n      \t.ExternalClass p,\n      \t.ExternalClass span,\n      \t.ExternalClass font,\n      \t.ExternalClass td,\n      \t.ExternalClass div {\n      \t\tline-height: 100%;\n      \t}\n\n      \t.apple-link a {\n      \t\tcolor: inherit !important;\n      \t\tfont-family: inherit !important;\n      \t\tfont-size: inherit !important;\n      \t\tfont-weight: inherit !important;\n      \t\tline-height: inherit !important;\n      \t\ttext-decoration: none !important;\n      \t}\n\n      \t#MessageViewBody a {\n      \t\tcolor: inherit;\n      \t\ttext-decoration: none;\n      \t\tfont-size: inherit;\n      \t\tfont-family: inherit;\n      \t\tfont-weight: inherit;\n      \t\tline-height: inherit;\n      \t}\n\n      \t.btn-primary table td:hover {\n      \t\tbackground-color: #34495e !important;\n      \t}\n\n      \t.btn-primary a:hover {\n      \t\tbackground-color: #34495e !important;\n      \t\tborder-color: #34495e !important;\n      \t}\n      }\n\n      /*-------------------------------------\n        MY CODE\n      --------------------------------------/*\n    </style></head><body><span class=\"preheader\">pagemail saved pages from ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(date.Format("02-01-2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 370, Col: 80}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"body\"><tr><td>&nbsp;</td><td class=\"container\"><div class=\"content\"><!-- START CENTERED WHITE CONTAINER --><table role=\"presentation\" class=\"main\"><!-- START MAIN CONTENT AREA --><tr><td class=\"wrapper\"><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td><p>Good morning ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 395, Col: 35}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><p>Here are all the pages you've saved since ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(date.Format("02-01-2006"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 397, Col: 83}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><ul class=\"main-list\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, page := range pages {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"main-list-item\"><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 templ.SafeURL = templ.URL(page.Url)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var6)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if page.Title != nil {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<b>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(*page.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 404, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</b> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if page.Description != nil {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("- ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var8 string
					templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(*page.Description)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 406, Col: 40}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<b>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(page.Url)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 409, Col: 31}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</b>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a> <i><span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var10 string
			templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(page.Created.Format("02/01 15:04"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `pagemail/internal/render/mail.templ`, Line: 412, Col: 61}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></i></li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></td></tr></table></td></tr><!-- END MAIN CONTENT AREA --></table><!-- END CENTERED WHITE CONTAINER --><!-- START FOOTER --><div class=\"footer\"><table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td class=\"content-block\"><span class=\"apple-link\">Thank's for subscribing to pagemail!</span><br>Don't like these emails? <a href=\"https://pagemail.io/account\">Unsubscribe</a>.</td></tr><tr><td class=\"content-block powered-by\">Powered by <a href=\"http://htmlemail.io\">HTMLemail</a>.</td></tr></table></div><!-- END FOOTER --></div></td><td>&nbsp;</td></tr></table></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
