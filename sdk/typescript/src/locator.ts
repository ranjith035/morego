export class Locator {
  constructor(
    private session: any,
    private strategy: string,
    private selector: string
  ) {}

  private async resolve(): Promise<string> {
    let stratEnum = 0;
    switch (this.strategy.toUpperCase()) {
      case 'ACCESSIBILITY_ID':
      case 'TEST_ID':
        stratEnum = 1;
        break;
      case 'RESOURCE_ID':
        stratEnum = 2;
        break;
      case 'TEXT':
        stratEnum = 3;
        break;
      case 'CLASS_NAME':
      case 'CLASS':
        stratEnum = 4;
        break;
      case 'XPATH':
        stratEnum = 5;
        break;
    }

    return new Promise((resolve, reject) => {
      this.session.device.driverClient.FindElement(
        {
          session_id: this.session.sessionId,
          locator: { strategy: stratEnum, selector: this.selector }
        },
        (err: any, response: any) => {
          if (err) return reject(err);
          if (!response.found) return reject(new Error(`Element not found: ${this.strategy}=${this.selector}`));
          resolve(response.element.element_id);
        }
      );
    });
  }

  async click(): Promise<void> {
    const elId = await this.resolve();
    return new Promise((resolve, reject) => {
      this.session.device.driverClient.Click(
        { driver_id: this.session.sessionId, element_id: elId },
        (err: any) => {
          if (err) return reject(err);
          resolve();
        }
      );
    });
  }

  async fill(value: string): Promise<void> {
    const elId = await this.resolve();
    return new Promise((resolve, reject) => {
      this.session.device.driverClient.Fill(
        { driver_id: this.session.sessionId, element_id: elId, value },
        (err: any) => {
          if (err) return reject(err);
          resolve();
        }
      );
    });
  }
}
