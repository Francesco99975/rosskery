import { LitElement, html, css } from "lit";
import { property, query } from "lit/decorators.js";

export class UCounter extends LitElement {
  @property({ type: String })
  name: string = "";

  @property({ type: Number })
  min: number = 1;

  @property({ type: Number })
  value: number = 1;

  @query("input") input!: HTMLInputElement;

  static styles = css`
    div.flex {
      display: flex;
      align-items: center;
      justify-content: space-evenly;
      width: 100%;
    }

    button {
      border-width: 4px;
      width: 3rem;
      height: 3rem;
      --tw-border-opacity: 1;
      border-color: rgb(var(--color-primary) / var(--tw-border-opacity));
      border-style: solid;
      border-radius: 9999px;
      text-align: center;
      --tw-text-opacity: 1;
      color: rgb(var(--color-primary) / var(--tw-text-opacity));
      font-weight: 700;
      font-size: 1.875rem;
      line-height: 2.25rem;
    }
  `;

  private _increase() {
    this.value = +this.value + 1;
    this._dispatchChangeEvent();
  }

  private _decrease() {
    if (+this.value > +this.min) {
      this.value = +this.value - 1;
    }
    this._dispatchChangeEvent();
  }

  private _dispatchChangeEvent() {
    this.dispatchEvent(
      new CustomEvent("input", {
        detail: { value: this.value },
        bubbles: true,
        composed: true,
      })
    );
  }

  connectedCallback() {
    super.connectedCallback();
    const form = this.closest("form");
    if (form) {
      form.addEventListener("reset", this._handleFormReset);
    }
  }

  disconnectedCallback() {
    const form = this.closest("form");
    if (form) {
      form.removeEventListener("reset", this._handleFormReset);
    }
    super.disconnectedCallback();
  }

  private _handleFormReset = () => {
    this.value = this.min;
  };

  get formValue() {
    return this.value;
  }

  protected render() {
    return html`
      <div class="flex">
        <input
          .id="${"qty" + this.name}"
          .name="${"qty" + this.name}"
          .min="${this.min.toString()}"
          .value="${this.value.toString()}"
          type="hidden"
        />
        <button
          type="button"
          id="${"dec" + this.name}"
          @click="${this._decrease}"
        >
          -
        </button>
        <span class="text-xl md:text-3xl font-bold p-2 text-center"
          >${this.value}</span
        >
        <button
          type="button"
          id="${"inc" + this.name}"
          @click="${this._increase}"
        >
          +
        </button>
      </div>
    `;
  }
}
