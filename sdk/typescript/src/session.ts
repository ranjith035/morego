import { Locator } from './locator';

export class Session {
  constructor(
    public device: any,
    public sessionId: string,
    public appId: string
  ) {}

  async close(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.device.client.CloseSession({ session_id: this.sessionId }, (err: any) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  async swipe(startX: number, startY: number, endX: number, endY: number, durationMs: number): Promise<void> {
    return new Promise((resolve, reject) => {
      this.device.driverClient.Swipe(
        {
          driver_id: this.sessionId,
          start: { x: startX, y: startY },
          end: { x: endX, y: endY },
          duration_ms: durationMs
        },
        (err: any) => {
          if (err) return reject(err);
          resolve();
        }
      );
    });
  }

  async getSource(format: string = 'xml'): Promise<string> {
    return new Promise((resolve, reject) => {
      this.device.driverClient.GetSource(
        { driver_id: this.sessionId, format },
        (err: any, response: any) => {
          if (err) return reject(err);
          resolve(response.source_data);
        }
      );
    });
  }

  locator(strategy: string, selector: string): Locator {
    return new Locator(this, strategy, selector);
  }
}
