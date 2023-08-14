import { retryPageLoading } from './retryPageLoading';

describe('retryPageLoading', () => {
  let windowSpy;
  let windowObj;

  const setUpWindow = (localStorageItem) => {
    windowObj = {
      localStorage: {
        getItem: jest.fn().mockImplementation(() => localStorageItem),
        setItem: jest.fn(),
      },
      location: {
        reload: jest.fn(),
      },
    };
    windowSpy.mockImplementation(() => windowObj);
  };

  beforeEach(() => {
    windowSpy = jest.spyOn(global, 'window', 'get');
  });

  afterEach(() => {
    windowSpy.mockRestore();
  });

  it('does not reload on non chuck errors', () => {
    setUpWindow(null);
    retryPageLoading({ name: 'SomethingError' });
    expect(windowObj.localStorage.getItem).not.toBeCalled();
    expect(windowObj.localStorage.setItem).not.toBeCalled();
    expect(windowObj.location.reload).not.toBeCalled();
  });

  it('reloads on first chuck error', () => {
    setUpWindow('false');
    retryPageLoading({ name: 'ChunkLoadError' });
    expect(windowObj.localStorage.getItem).toBeCalled();
    expect(windowObj.localStorage.setItem).toBeCalledWith(expect.any(String), 'true');
    expect(windowObj.location.reload).toBeCalled();
  });

  it('does not reload on 2nd chuck error', () => {
    setUpWindow('true');
    retryPageLoading({ name: 'ChunkLoadError' });
    expect(windowObj.localStorage.getItem).toBeCalled();
    expect(windowObj.localStorage.setItem).toBeCalledWith(expect.any(String), 'false');
    expect(windowObj.location.reload).not.toBeCalled();
  });
});
