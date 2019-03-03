import { VoicesModule } from './voices.module';

describe('VoicesModule', () => {
  let voicesModule: VoicesModule;

  beforeEach(() => {
    voicesModule = new VoicesModule();
  });

  it('should create an instance', () => {
    expect(voicesModule).toBeTruthy();
  });
});
