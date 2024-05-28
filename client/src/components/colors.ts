import { LitElement, html, css } from "lit";
import { property } from "lit/decorators.js";

export class ColorPicker extends LitElement {
  @property({ type: String }) value: string = "#000000"; // Default value

  static styles = css`
    :host {
      display: inline-block;
    }
    input {
      width: 100px;
      height: 30px;
      padding: 5px;
      border: 1px solid #ccc;
    }
  `;

  render() {
    return html`
      <input type="color" .value="${this.value}" @input="${this.onChange}" />
    `;
  }

  onChange(event: Event) {
    const inputElement = event.target as HTMLInputElement;
    this.value = inputElement.value;
    this.dispatchEvent(
      new CustomEvent("input", { detail: { value: this.value } })
    );
  }

  get form(): HTMLFormElement | null {
    return this.closest("form");
  }

  get name(): string | null {
    return this.getAttribute("name");
  }

  get validity(): ValidityState {
    return this.inputElement.validity;
  }

  get validationMessage(): string {
    return this.inputElement.validationMessage;
  }

  checkValidity(): boolean {
    return this.inputElement.checkValidity();
  }

  reportValidity(): boolean {
    return this.inputElement.reportValidity();
  }

  submit(): void {
    if (this.form) {
      this.form.submit();
    }
  }

  get inputElement(): HTMLInputElement {
    return this.shadowRoot!.querySelector("input") as HTMLInputElement;
  }

  updated(changedProperties: Map<string | number | symbol, unknown>): void {
    if (changedProperties.has("value")) {
      this.inputElement.value = this.value;
    }
  }
}
