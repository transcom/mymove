import handleScroll from './handleScroll';

const sections = [{ id: '1' }, { id: '2' }];

const setActiveSection = jest.fn();

describe('handleScroll', () => {
  it('calls setActiveSection with expected new section', () => {
    jest.spyOn(document, 'querySelector').mockImplementation((selector) => {
      return selector;
    });
    handleScroll(sections, '2', setActiveSection);
    expect(setActiveSection).toHaveBeenCalled();
  });
});
