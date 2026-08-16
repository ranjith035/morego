export class Locator {
  private indexVal: number = 0;
  private parentVal?: Locator;
  private constraintsVal: Array<{ direction: string; target: Locator; distancePx?: number }> = [];

  constructor(
    private session: any,
    private strategy: string,
    private selector: string
  ) {}

  private toProto(): any {
    let strat = 'LOCATOR_STRATEGY_UNSPECIFIED';
    switch (this.strategy.toUpperCase()) {
      case 'ACCESSIBILITY_ID':
        strat = 'LOCATOR_STRATEGY_ACCESSIBILITY_ID';
        break;
      case 'TEST_ID':
        strat = 'LOCATOR_STRATEGY_TEST_ID';
        break;
      case 'ROLE':
        strat = 'LOCATOR_STRATEGY_ROLE';
        break;
      case 'TEXT':
        strat = 'LOCATOR_STRATEGY_TEXT';
        break;
      case 'PLACEHOLDER':
        strat = 'LOCATOR_STRATEGY_PLACEHOLDER';
        break;
      case 'LABEL':
        strat = 'LOCATOR_STRATEGY_LABEL';
        break;
      case 'RESOURCE_ID':
        strat = 'LOCATOR_STRATEGY_RESOURCE_ID';
        break;
      case 'XPATH':
        strat = 'LOCATOR_STRATEGY_XPATH';
        break;
    }

    const p: any = {
      strategy: strat,
      selector: this.selector,
      index: this.indexVal
    };

    if (this.parentVal) {
      p.parent = this.parentVal.toProto();
    }

    if (this.constraintsVal && this.constraintsVal.length > 0) {
      p.constraints = this.constraintsVal.map((c: any) => {
        let dir = 'RELATIVE_DIRECTION_UNSPECIFIED';
        switch (c.direction.toUpperCase()) {
          case 'ABOVE':
            dir = 'RELATIVE_DIRECTION_ABOVE';
            break;
          case 'BELOW':
            dir = 'RELATIVE_DIRECTION_BELOW';
            break;
          case 'LEFT_OF':
            dir = 'RELATIVE_DIRECTION_LEFT_OF';
            break;
          case 'RIGHT_OF':
            dir = 'RELATIVE_DIRECTION_RIGHT_OF';
            break;
        }
        return {
          direction: dir,
          target: c.target.toProto(),
          distance_px: c.distancePx || 0
        };
      });
    }

    return p;
  }

  private async resolve(): Promise<string> {
    return new Promise((resolve, reject) => {
      this.session.device.driverClient.FindElement(
        {
          session_id: this.session.sessionId,
          locator: this.toProto()
        },
        (err: any, response: any) => {
          if (err) return reject(err);
          if (!response.found) return reject(new Error(`Element not found: ${this.strategy}=${this.selector}`));
          resolve(response.element.element_id);
        }
      );
    });
  }

  nth(index: number): Locator {
    const loc = new Locator(this.session, this.strategy, this.selector);
    loc.indexVal = index;
    loc.parentVal = this.parentVal;
    loc.constraintsVal = [...this.constraintsVal];
    return loc;
  }

  locator(strategy: string, selector: string): Locator {
    const loc = new Locator(this.session, strategy, selector);
    loc.parentVal = this;
    return loc;
  }

  getByText(text: string): Locator {
    return this.locator('TEXT', text);
  }

  getByRole(role: string): Locator {
    return this.locator('ROLE', role);
  }

  getByLabel(label: string): Locator {
    return this.locator('LABEL', label);
  }

  getByPlaceholder(placeholder: string): Locator {
    return this.locator('PLACEHOLDER', placeholder);
  }

  getByAccessibilityID(id: string): Locator {
    return this.locator('ACCESSIBILITY_ID', id);
  }

  getByTestID(id: string): Locator {
    return this.locator('TEST_ID', id);
  }

  above(other: Locator, distancePx?: number): Locator {
    this.constraintsVal.push({ direction: 'ABOVE', target: other, distancePx });
    return this;
  }

  below(other: Locator, distancePx?: number): Locator {
    this.constraintsVal.push({ direction: 'BELOW', target: other, distancePx });
    return this;
  }

  leftOf(other: Locator, distancePx?: number): Locator {
    this.constraintsVal.push({ direction: 'LEFT_OF', target: other, distancePx });
    return this;
  }

  rightOf(other: Locator, distancePx?: number): Locator {
    this.constraintsVal.push({ direction: 'RIGHT_OF', target: other, distancePx });
    return this;
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
