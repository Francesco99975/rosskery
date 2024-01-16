export class UCounter extends HTMLElement {
  idd: string = "";
  min: number = 0;
  value: number = 0;

  constructor() {
    super();

    this.idd = this.getAttribute("idd")!;

    this.min = +this.getAttribute("min")!;

    this.value = +this.getAttribute("value")!;
  }

  connectedCallback() {
    this.render();
  }

  increase() {
    this.value = +this.value + 1;
    console.log("inc", this.value);
    this.render();
  }

  decrease() {
    if (+this.value > +this.min) {
      this.value = +this.value - 1;
      console.log("dec", this.value);
      this.render();
    }
  }

  render() {
    this.innerHTML = `
      <div class="flex items-center justify-evenly w-full">
        <input
          id="${"qty" + this.idd}"
          min="${this.min.toString()}"
          value="${this.value.toString()}"
          type="hidden"
        />
        <button id="${
          "dec" + this.idd
        }" class="border-4 w-12 h-12 border-primary border-solid rounded-full text-center text-primary font-bold text-3xl">-</button>
        <span class="text-xl md:text-3xl font-bold p-2 text-center">${
          this.value
        }</span>
        <button id="${
          "inc" + this.idd
        }" class="border-4 w-12 h-12 border-primary border-solid rounded-full text-center text-primary font-bold text-3xl">+</button>
      </div>
    `;
    document
      ?.getElementById("dec" + this.idd)
      ?.addEventListener("click", this.decrease.bind(this));
    document
      ?.getElementById("inc" + this.idd)
      ?.addEventListener("click", this.increase.bind(this));
  }
}
