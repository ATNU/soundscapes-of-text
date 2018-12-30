import { TextModule } from './text.module';

describe('TextModule', () => {
  let textModule: TextModule;

  beforeEach(() => {
    textModule = new TextModule();
  });

  it('should create an instance', () => {
    expect(textModule).toBeTruthy();
  });
});
