import React from 'react';
import { shallow } from 'enzyme';
import { WizardPage } from 'shared/WizardPage';
describe('given a WizardPage', () => {
  let wrapper, buttons;
  const pageList = ['1', '2', '3'];
  const submit = jest.fn();
  const mockPush = jest.fn();
  describe('when handler is not async', () => {
    describe('when there is a pageIsValid prop set', () => {
      describe('when pageIsValid is false', () => {
        describe('when on the first page', () => {
          beforeEach(() => {
            const continueToNextPage = false;

            wrapper = shallow(
              <WizardPage
                handleSubmit={submit}
                pageList={pageList}
                pageKey="1"
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
            const pageIsValid = false;

            wrapper = shallow(
              <WizardPage
                handleSubmit={submit}
                pageList={pageList}
                pageKey="3"
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
          const continueToNextPage = true;

          wrapper = shallow(
            <WizardPage
              handleSubmit={submit}
              pageList={pageList}
              pageKey="1"
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
    describe('when there is an error', () => {
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
              isAsync={true}
              hasSucceeded={false}
              error={{ message: 'Something bad happened' }}
            >
              <div>This is page 2</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('it shows an error alert before its child', () => {
          const childContainer = wrapper.find('div.usa-width-one-whole');
          expect(childContainer.length).toBe(3);
          expect(childContainer.first().text()).toBe('<Alert />');
        });
        it('it renders button for cancel, back, next', () => {
          expect(buttons.length).toBe(3);
        });
        it('the cancel button is first and is enabled', () => {
          const cancelButton = buttons.first();
          expect(cancelButton.text()).toBe('Cancel');
          expect(cancelButton.prop('disabled')).toBe(false);
        });
        it('the back button is second and is enabled', () => {
          const backButton = buttons.at(1);
          expect(backButton.text()).toBe('Back');
          expect(backButton.prop('disabled')).toBe(false);
        });
        it('the next button is last and is enabled', () => {
          const nextButton = buttons.last();
          expect(nextButton.text()).toBe('Next');
          expect(nextButton.prop('disabled')).toBe(false);
        });
      });
    });
    describe('when page is not dirty', () => {
      beforeEach(() => {
        mockPush.mockClear();
        wrapper = shallow(
          <WizardPage
            handleSubmit={submit}
            pageList={pageList}
            pageKey="2"
            push={mockPush}
            match={{}}
            isAsync={true}
            hasSucceeded={false}
            pageIsDirty={false}
          >
            <div>This is page 2</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      it('the previous button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('1');
        });
        it('submit is not called', () => {
          expect(submit.mock.calls.length).toBe(0);
        });
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
        it('submit is not called', () => {
          expect(submit.mock.calls.length).toBe(0);
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
      afterEach(() => mockPush.mockClear());
      it('it starts on the first page', () => {
        const childContainer = wrapper.find('div.usa-width-one-whole');
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 1');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is disabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(true);
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();

        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 2');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the back button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('1');
        });
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
      afterEach(() => {
        submit.mockClear();
      });

      it('it shows its child', () => {
        const childContainer = wrapper.find('div.usa-width-one-whole');
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 3');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the back button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('2');
        });
      });
      it('the cancel button is first and is enabled', () => {
        const saveButton = buttons.first();
        expect(saveButton.text()).toBe('Cancel');
        expect(saveButton.prop('disabled')).toBe(false);
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
  describe('when handler is async', () => {
    describe('when there is a pageIsValid prop set', () => {
      describe('when pageIsValid is false', () => {
        describe('when on the first page', () => {
          beforeEach(() => {
            const continueToNextPage = false;

            wrapper = shallow(
              <WizardPage
                handleSubmit={submit}
                pageList={pageList}
                pageKey="1"
                pageIsValid={continueToNextPage}
                match={{}}
                isAsync={true}
                hasSucceeded={false}
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
            const pageIsValid = false;

            wrapper = shallow(
              <WizardPage
                handleSubmit={submit}
                pageList={pageList}
                pageKey="3"
                pageIsValid={pageIsValid}
                match={{}}
                isAsync={true}
                hasSucceeded={false}
              >
                <div>This is page 1</div>
              </WizardPage>,
            );
            buttons = wrapper.find('button');
          });
          it('the back button is second and is disabled', () => {
            const prevButton = buttons.at(1);
            expect(prevButton.text()).toBe('Back');
            expect(prevButton.prop('disabled')).toBe(true);
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
          const continueToNextPage = true;

          wrapper = shallow(
            <WizardPage
              handleSubmit={submit}
              pageList={pageList}
              pageKey="1"
              pageIsValid={continueToNextPage}
              isAsync={true}
              hasSucceeded={false}
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
    describe('when there is an error', () => {
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
              isAsync={true}
              hasSucceeded={false}
              error={{ message: 'Something bad happened' }}
            >
              <div>This is page 2</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('it shows an error alert before its child', () => {
          const childContainer = wrapper.find('div.usa-width-one-whole');
          expect(childContainer.length).toBe(3);
          expect(childContainer.first().text()).toBe('<Alert />');
        });
        it('it renders button for cancel, back, next', () => {
          expect(buttons.length).toBe(3);
        });
        it('the back button is second and is enabled', () => {
          const prevButton = buttons.at(1);
          expect(prevButton.text()).toBe('Back');
          expect(prevButton.prop('disabled')).toBe(false);
        });
        it('the cancel button is first and is enabled', () => {
          const prevButton = buttons.first();
          expect(prevButton.text()).toBe('Cancel');
          expect(prevButton.prop('disabled')).toBe(false);
        });
        it('the next button is last and is enabled', () => {
          const nextButton = buttons.last();
          expect(nextButton.text()).toBe('Next');
          expect(nextButton.prop('disabled')).toBe(false);
        });
      });
    });
    describe('when page is not dirty', () => {
      beforeEach(() => {
        mockPush.mockClear();
        wrapper = shallow(
          <WizardPage
            handleSubmit={submit}
            pageList={pageList}
            pageKey="2"
            push={mockPush}
            match={{}}
            isAsync={true}
            hasSucceeded={false}
            pageIsDirty={false}
          >
            <div>This is page 2</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => {
        submit.mockClear();
      });
      it('the back button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('submit is not called', () => {
          expect(submit.mock.calls.length).toBe(0);
        });
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
        it('submit is not called', () => {
          expect(submit.mock.calls.length).toBe(0);
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
            isAsync={true}
            hasSucceeded={false}
          >
            <div>This is page 1</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => {
        submit.mockClear();
      });
      it('it starts on the first page', () => {
        const childContainer = wrapper.find('div.usa-width-one-whole');
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 1');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is disabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(true);
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
        it('transitionFunc is set to getNextPage', () => {
          const state = wrapper.state();
          expect(state.transitionFunc.name).toBe('getNextPagePath');
        });
        it('submit is called', () => {
          expect(submit.mock.calls.length).toBe(1);
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
            isAsync={true}
            hasSucceeded={false}
          >
            <div>This is page 2</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => {
        submit.mockClear();
      });
      it('it shows its child', () => {
        const childContainer = wrapper.find('div.usa-width-one-whole');
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 2');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('transitionFunc is set to getPrevPage', () => {
          const state = wrapper.state();
          expect(state.transitionFunc.name).toBe('getPreviousPagePath');
        });
        it('submit is called', () => {
          expect(submit.mock.calls.length).toBe(1);
        });
      });
      it('the cancel button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
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
        it('transitionFunc is set to getNextPage', () => {
          const state = wrapper.state();
          expect(state.transitionFunc.name).toBe('getNextPagePath');
        });

        it('submit is called', () => {
          expect(submit.mock.calls.length).toBe(1);
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
            isAsync={true}
            hasSucceeded={false}
          >
            <div>This is page 3</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => {
        submit.mockClear();
      });

      it('it shows its child', () => {
        const childContainer = wrapper.find('div.usa-width-one-whole');
        expect(childContainer.length).toBe(2);
        expect(childContainer.first().text()).toBe('This is page 3');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is second and is enabled', () => {
        const prevButton = buttons.at(1);
        expect(prevButton.text()).toBe('Back');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = buttons.at(1);
          prevButton.simulate('click');
        });
        it('submit is called', () => {
          expect(submit.mock.calls.length).toBe(1);
        });
      });
      it('the cancel button is first and is enabled', () => {
        const saveButton = buttons.first();
        expect(saveButton.text()).toBe('Cancel');
        expect(saveButton.prop('disabled')).toBe(false);
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
  describe('when there is an additionalParams prop', () => {
    beforeEach(() => {
      mockPush.mockClear();
      wrapper = shallow(
        <WizardPage
          pageList={['page1', 'anotherPage/:foo/:bar']}
          pageKey="page1"
          match={{ params: { foo: 'dvorak' } }}
          push={mockPush}
          handleSubmit={() => undefined}
          additionalParams={{ bar: 'querty' }}
        >
          <div>This is page 1</div>
        </WizardPage>,
      );
      buttons = wrapper.find('button');
    });
    describe('when the next button is clicked', () => {
      beforeEach(() => {
        const nextButton = buttons.last();
        nextButton.simulate('click');
      });
      it('push gets a page with the additionalParams expanded', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('anotherPage/dvorak/querty');
      });
    });
  });
});
