import React from 'react';
import { mount, shallow } from 'enzyme';
import { WizardPage } from 'shared/WizardPage';
describe('given a WizardPage', () => {
  let wrapper, buttons;
  const pageList = ['1', '2', '3'];
  const submit = jest.fn();
  const mockPush = jest.fn();
  describe('when there is a pageIsValid prop set', () => {
    describe('when pageIsValid is false', () => {
      describe('when on the first page', () => {
        beforeEach(() => {
          history = [];
          const continueToNextPage = false;

          wrapper = shallow(
            <WizardPage
              handleSubmit={submit}
              pageList={pageList}
              pageKey="1"
              history={history}
              pageIsValid={continueToNextPage}
              match={{}}
            >
              <div>This is page 1</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('the next button is last and is disabled', () => {
          const nextButton = buttons.last();
          expect(nextButton.text()).toBe('Next');
          expect(nextButton.prop('disabled')).toBeTruthy();
        });
      });
      describe('when on the last page', () => {
        beforeEach(() => {
          history = [];
          const pageIsValid = false;

          wrapper = shallow(
            <WizardPage
              handleSubmit={submit}
              pageList={pageList}
              pageKey="3"
              history={history}
              pageIsValid={pageIsValid}
              match={{}}
            >
              <div>This is page 1</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('the complete button is last and is disabled', () => {
          const nextButton = buttons.last();
          expect(nextButton.text()).toBe('Complete');
          expect(nextButton.prop('disabled')).toBeTruthy();
        });
      });
    });
    describe('when pageIsValid is true', () => {
      beforeEach(() => {
        history = [];
        const continueToNextPage = true;

        wrapper = shallow(
          <WizardPage
            handleSubmit={submit}
            pageList={pageList}
            pageKey="1"
            history={history}
            pageIsValid={continueToNextPage}
          >
            <div>This is page 1</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      it('the next button is last and is enabled', () => {
        const nextButton = buttons.last();
        expect(nextButton.text()).toBe('Next');
        expect(nextButton.prop('disabled')).toBeFalsy();
      });
    });
  });
  describe('when on the first page', () => {
    beforeEach(() => {
      wrapper = shallow(
        <WizardPage
          handleSubmit={submit}
          pageList={pageList}
          pageKey="1"
          push={mockPush}
          match={{}}
        >
          <div>This is page 1</div>
        </WizardPage>,
      );
      buttons = wrapper.find('button');
    });
    it('it starts on the first page', () => {
      const childContainer = wrapper.find('div.usa-width-one-whole');
      expect(childContainer.length).toBe(1);
      expect(childContainer.first().text()).toBe('This is page 1');
    });
    it('it renders button for prev, save, next', () => {
      expect(buttons.length).toBe(3);
    });
    it('the previous button is first and is disabled', () => {
      const prevButton = buttons.first();
      expect(prevButton.text()).toBe('Prev');
      expect(prevButton.prop('disabled')).toBe(true);
    });
    it('the save for later button is second and is disabled', () => {
      const prevButton = buttons.at(1);
      expect(prevButton.text()).toBe('Save for later');
      expect(prevButton.prop('disabled')).toBe(true);
    });
    it('the next button is last and is enabled', () => {
      const nextButton = buttons.last();
      expect(nextButton.text()).toBe('Next');
      expect(nextButton.prop('disabled')).toBeFalsy();
    });

    describe('when the next button is clicked', () => {
      beforeEach(() => {
        const nextButton = buttons.last();
        nextButton.simulate('click');
      });
      it('push gets the next page', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('2');
      });
    });
  });

  describe('when on the middle page', () => {
    beforeEach(() => {
      mockPush.mockClear();
      wrapper = shallow(
        <WizardPage
          handleSubmit={submit}
          pageList={pageList}
          pageKey="2"
          push={mockPush}
          match={{}}
        >
          <div>This is page 2</div>
        </WizardPage>,
      );
      buttons = wrapper.find('button');
    });
    it('it shows its child', () => {
      const childContainer = wrapper.find('div.usa-width-one-whole');
      expect(childContainer.length).toBe(1);
      expect(childContainer.first().text()).toBe('This is page 2');
    });
    it('it renders button for prev, save, next', () => {
      expect(buttons.length).toBe(3);
    });
    it('the previous button is first and is enabled', () => {
      const prevButton = buttons.first();
      expect(prevButton.text()).toBe('Prev');
      expect(prevButton.prop('disabled')).toBe(false);
    });
    describe('when the prev button is clicked', () => {
      beforeEach(() => {
        const prevButton = buttons.first();
        prevButton.simulate('click');
      });
      it('push gets the prev page', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('1');
      });
    });
    it('the save for later button is second and is disabled', () => {
      const prevButton = buttons.at(1);
      expect(prevButton.text()).toBe('Save for later');
      expect(prevButton.prop('disabled')).toBe(true);
    });
    it('the next button is last and is enabled', () => {
      const nextButton = buttons.last();
      expect(nextButton.text()).toBe('Next');
      expect(nextButton.prop('disabled')).toBeFalsy();
    });
    describe('when the next button is clicked', () => {
      beforeEach(() => {
        const nextButton = buttons.last();
        nextButton.simulate('click');
      });
      it('push gets the next page', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('3');
      });
    });
  });
  describe('when on the last page', () => {
    beforeEach(() => {
      mockPush.mockClear();
      wrapper = shallow(
        <WizardPage
          handleSubmit={submit}
          pageList={pageList}
          pageKey="3"
          push={mockPush}
          match={{}}
        >
          <div>This is page 3</div>
        </WizardPage>,
      );
      buttons = wrapper.find('button');
    });
    it('it shows its child', () => {
      const childContainer = wrapper.find('div.usa-width-one-whole');
      expect(childContainer.length).toBe(1);
      expect(childContainer.first().text()).toBe('This is page 3');
    });
    it('it renders button for prev, save, next', () => {
      expect(buttons.length).toBe(3);
    });
    it('the previous button is first and is enabled', () => {
      const prevButton = buttons.first();
      expect(prevButton.text()).toBe('Prev');
      expect(prevButton.prop('disabled')).toBe(false);
    });
    describe('when the prev button is clicked', () => {
      beforeEach(() => {
        const prevButton = buttons.first();
        prevButton.simulate('click');
      });
      it('push gets the prev page', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('2');
      });
    });
    it('the save for later button is second and is disabled', () => {
      const saveButton = buttons.at(1);
      expect(saveButton.text()).toBe('Save for later');
      expect(saveButton.prop('disabled')).toBe(true);
    });
    it('the Complete button is last and is enabled', () => {
      const nextButton = buttons.last();
      expect(nextButton.text()).toBe('Complete');
      expect(nextButton.prop('disabled')).toBeFalsy();
    });
    describe('when the complete button is clicked', () => {
      beforeEach(() => {
        const nextButton = buttons.last();
        nextButton.simulate('click');
      });
      it('submit is called', () => {
        expect(submit.mock.calls.length).toBe(1);
      });
    });
  });
});
