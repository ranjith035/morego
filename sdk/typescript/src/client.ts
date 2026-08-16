import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import * as path from 'path';
import { Session } from './session';

export class Device {
  private client: any;
  public driverClient: any;
  public aiClient: any;

  constructor(private address: string, private deviceId: string) {
    const protoPath = path.resolve(__dirname, '../../../proto');
    const packageDefinition = protoLoader.loadSync(
      [
        path.join(protoPath, 'session.proto'),
        path.join(protoPath, 'driver.proto'),
        path.join(protoPath, 'ai.proto'),
        path.join(protoPath, 'locator.proto')
      ],
      {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true,
        includeDirs: [path.resolve(protoPath, '..'), protoPath]
      }
    );
    const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
    const automation = protoDescriptor.automation.v1;

    this.client = new automation.SessionService(address, grpc.credentials.createInsecure());
    this.driverClient = new automation.DriverService(address, grpc.credentials.createInsecure());
    this.aiClient = new automation.AIService(address, grpc.credentials.createInsecure());
  }

  static async connect(address: string, deviceId: string): Promise<Device> {
    return new Device(address, deviceId);
  }

  async newSession(appId: string, capabilities: Record<string, string> = {}): Promise<Session> {
    return new Promise((resolve, reject) => {
      this.client.CreateSession(
        { device_id: this.deviceId, app_id: appId, capabilities },
        (err: any, response: any) => {
          if (err) return reject(err);
          resolve(new Session(this, response.session_id, appId));
        }
      );
    });
  }

  close() {
    this.client.close();
    this.driverClient.close();
    this.aiClient.close();
  }
}
