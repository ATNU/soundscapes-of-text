import { EncodingModule } from './encoding.module';

describe('EncodingModule', () => {
  let encodingModule: EncodingModule;

  beforeEach(() => {
    encodingModule = new EncodingModule();
  });

  it('should create an instance', () => {
    expect(encodingModule).toBeTruthy();
  });
});
