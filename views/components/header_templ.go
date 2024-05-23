// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/Francesco99975/rosskery/views/icons"

func Header(message string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<header hx-boost=\"true\" class=\"grid grid-cols-3 gap-2 place-items-center bg-std text-center text-primary w-full h-24 p-4 sticky top-0 right-0 z-20 shadow-md border-b-2 border-b-primary rounded-b-lg\"><nav class=\"md:w-auto\"><!--")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := ` Burger menu icon for small screens `
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("--><div id=\"burgerMenu\" class=\"burger-menu md:hidden cursor-pointer\"><div id=\"bar1\" class=\"bar w-6 h-1 bg-primary my-1 rounded transition-transform transform rotate-0\"></div><div id=\"bar2\" class=\"bar w-6 h-1 bg-primary my-1 rounded transition-transform transform rotate-0\"></div><div id=\"bar3\" class=\"bar w-6 h-1 bg-primary my-1 rounded transition-transform transform rotate-0\"></div></div><!--")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := ` Navigation links for larger screens `
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("--><ul id=\"navLinks\" class=\"nav-links md:flex flex-row space-x-4 hidden\"><li><a href=\"/shop\" class=\"text-primary text-lg md:text-xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var4 := `Shop`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></li><li><a href=\"/gallery\" class=\"text-primary text-lg md:text-xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var5 := `Gallery`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var5)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></li></ul><!--")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var6 := ` Navigation links for mobile view `
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var6)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("--><ul id=\"mobileNavLinks\" class=\"nav-links-mobile md:hidden absolute top-24 left-0 w-full hidden z-30 transition-all ease-in\"><li class=\"bg-std w-full px-4 py-2\"><a href=\"/shop\" class=\"text-primary text-center text-xl md:text-2xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var7 := `Shop`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var7)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></li><li class=\"bg-std w-full px-4 py-2\"><a href=\"/gallery\" class=\"text-primary text-center text-xl md:text-2xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var8 := `Gallery`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var8)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></li></ul></nav><div class=\"flex items-center p-2\"><h1 class=\"text-3xl\"><a href=\"/\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var9 := `Rosskery`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var9)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></h1></div><a href=\"/checkout\" class=\"flex justify-center items-center relative\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = icons.BagIcon().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div hx-get=\"/bag\" hx-trigger=\"load\" hx-swap=\"outerHTML\"></div></a> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(message) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"mbg\" class=\"w-full bg-primary whitespace-nowrap overflow-hidden flex justify-center items-center fixed top-24 left-0 z-30\"><span id=\"mtx\" class=\"p-1 text-sm font-bold text-std animate-pacman\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var10 string = message
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var11 := `

            var burgerMenu = document.getElementById('burgerMenu');
            var navLinks = document.getElementById('mobileNavLinks');
            var bar1 = document.getElementById('bar1');
            var bar2 = document.getElementById('bar2');
            var bar3 = document.getElementById('bar3');

            burgerMenu.addEventListener('click', function () {
                navLinks.classList.toggle('hidden');
                  if (bar1.classList.contains('rotate-0')) {
                    bar1.classList.remove('rotate-0');
                    bar1.classList.add('rotate-45', 'translate-y-2');

                    bar2.classList.remove('rotate-0');
                    bar2.classList.add('opacity-0');

                    bar3.classList.remove('rotate-0');
                    bar3.classList.add('-rotate-45', '-translate-y-2');
                } else {
                    bar1.classList.remove('rotate-45', 'translate-y-2');
                    bar1.classList.add('rotate-0');

                    bar2.classList.remove('opacity-0');
                    bar3.classList.remove('-rotate-45', '-translate-y-2');
                    bar3.classList.add('rotate-0');
                }
            });

						 // Adjusting span width to wrap the text continuously
						const span = document.getElementById('mtx');
						const div = document.getElementById('mbg');
						if(div && span) {
							const divWidth = div.offsetWidth;
							const spanWidth = span.offsetWidth;
							const clonesNeeded = Math.ceil(divWidth / spanWidth) + 1;

							for (let i = 0; i < clonesNeeded; i++) {
									const clone = span.cloneNode(true);
									span.parentNode.appendChild(clone);
							}
						}

        `
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var11)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</script></header>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
