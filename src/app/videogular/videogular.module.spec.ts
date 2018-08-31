import { VideogularModule } from './videogular.module';

describe('VideogularModule', () => {
  let videogularModule: VideogularModule;

  beforeEach(() => {
    videogularModule = new VideogularModule();
  });

  it('should create an instance', () => {
    expect(videogularModule).toBeTruthy();
  });
});
