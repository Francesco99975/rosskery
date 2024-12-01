// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/views/layouts"
)

func Checkout(site models.Site, cartPreview *models.CartPreview, overbookedData string, csrf string, nonce string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main class=\"flex flex-col gap-2 w-full bg-primary min-h-screen\"><div class=\"w-[90%] md:max-w-7xl mx-auto bg-std p-4 md:p-6 rounded-lg shadow-md mt-3\"><section class=\"mb-6 text-primary\"><h2 class=\"text-2xl md:text-3xl font-bold mb-4\">Your Bag</h2><div id=\"cart-items\" class=\"space-y-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, item := range cartPreview.Items {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col md:flex-row justify-between items-start md:items-center border-b-2 border-primary pb-4\"><div><h3 class=\"text-lg font-semibold\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(helpers.Capitalize(item.Product.Name))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 20, Col: 82}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h3>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if !item.Product.Weighed {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p>Quantity: ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var4 string
					templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(item.Quantity))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 22, Col: 50}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<p>Quantity: ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(float64(item.Quantity) / 10))
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 24, Col: 62}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("lb</p>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"mt-2 md:mt-0 text-right md:text-left\"><p class=\"text-lg\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(helpers.FormatPrice(float64(item.Subtotal) / 100.0))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 28, Col: 81}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"mt-6\"><div class=\"flex justify-between text-xl font-bold mt-4 text-accent\"><p>Total:</p><p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(helpers.FormatPrice(float64(cartPreview.Total) / 100.0))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 36, Col: 67}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div></section><!-- Customer Information Form Section --><section class=\"mb-6 text-primary\"><h2 class=\"text-xl md:text-2xl font-bold mb-4\">Customer Information</h2><form id=\"checkout-form\" hx-post=\"/orders\" id=\"checkout-form\" class=\"space-y-4\" hx-target=\"body\" hx-boost=\"true\"><input type=\"hidden\" name=\"dd\" id=\"dd\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(overbookedData)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 44, Col: 67}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" name=\"_csrf\" id=\"_csrf\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(csrf)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `checkout.templ`, Line: 45, Col: 63}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"grid grid-cols-1 md:grid-cols-2 gap-4\"><div><label for=\"email\" class=\"block text-sm font-medium\">Email</label> <input type=\"email\" id=\"email\" name=\"email\" required class=\"mt-1 block w-full rounded-md border-primary shadow-sm focus:ring-iaccent focus:border-accent p-1\"></div><div><label for=\"fullname\" class=\"block text-sm font-medium\">Full Name</label> <input type=\"text\" id=\"fullname\" name=\"fullname\" required class=\"mt-1 block w-full rounded-md border-primaryshadow-sm focus:ring-accent focus:border-accent p-1\"></div><div class=\"md:col-span-2\"><label for=\"address\" class=\"block text-sm font-medium\">Address</label> <input type=\"text\" id=\"address\" name=\"address\" required hx-get=\"/address\" hx-trigger=\"keyup changed delay:500ms\" hx-target=\"#suggestions\" autocomplete=\"off\" class=\"mt-1 block w-full rounded-md border-primaryshadow-sm focus:ring-accent focus:border-accent p-1\"><div id=\"suggestions\" class=\"border border-gray-300 mt-2 rounded bg-white shadow-lg\"></div></div><div><label for=\"phone\" class=\"block text-sm font-medium\">Phone Number</label> <input type=\"tel\" id=\"phone\" name=\"phone\" required class=\"mt-1 block w-full rounded-md border-primaryshadow-sm focus:ring-accent focus:border-accent p-1\"></div><div><label for=\"pickuptime\" class=\"block text-sm font-medium\">Pickup Time</label> <input type=\"hidden\" id=\"pickuptime\" name=\"pickuptime\" required class=\"mt-1 block w-full rounded-md border-primaryshadow-sm focus:ring-accent focus:border-accent p-1\"></div></div><!-- Payment Method Section --><section><h2 class=\"text-xl md:text-2xl font-bold mb-4\">Payment Method</h2><div class=\"flex space-x-2 border-[3px] border-accent rounded-xl select-none md:w-1/3\"><label class=\"radio flex flex-grow items-center justify-center rounded-lg p-1 cursor-pointer\"><input type=\"radio\" name=\"method\" value=\"stripe\" class=\"peer hidden\" checked=\"\"> <span class=\"tracking-widest peer-checked:bg-primary peer-checked:text-std text-primary p-2 rounded-lg transition duration-150 ease-in-out\">Pay Online</span></label> <label class=\"radio flex flex-grow items-center justify-center rounded-lg p-1 cursor-pointer\"><input type=\"radio\" name=\"method\" value=\"cash\" class=\"peer hidden\"> <span class=\"tracking-widest peer-checked:bg-primary peer-checked:text-std text-primary p-2 rounded-lg transition duration-150 ease-in-out\">Cash at Pickup</span></label></div></section><button type=\"submit\" form=\"checkout-form\" class=\"mt-6 w-full bg-primary text-std py-3 rounded-lg font-bold text-lg hover:bg-accent\">Place Order</button><div id=\"errors\"></div></form></section></div></main>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Payment(site, nonce, []string{"assets/dist/checkout.css"}, nil, []string{"/assets/dist/checkout.js"}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
