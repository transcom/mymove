import React from 'react';
import { shallow } from 'enzyme';
import { WizardPage } from 'shared/WizardPage';
describe('given a WizardPage', () => {
  let wrapper, buttons;
  const submit = jest.fn();
  const mockPush = jest.fn();

  const minProps = {
    handleSubmit: jest.fn(),
    pageList: ['1', '2', '3'],
    pageKey: '1',
  };
  describe('Component renders', () => {
    expect(shallow(<WizardPage {...minProps} />).length).toEqual(1);
  });
  describe('when handler is not async', () => {
    describe('when there is a pageIsValid prop set', () => {
      describe('when pageIsValid is false', () => {
        describe('when on the first page', () => {
          beforeEach(() => {
            const continueToNextPage = false;

            wrapper = shallow(
              <WizardPage {...minProps} pageIsValid={continueToNextPage} match={{}}>
                <div>This is page 1</div>
              </WizardPage>,
            );
            buttons = wrapper.find('button');
          });
          it('the next button is first and is disabled', () => {
            const nextButton = buttons.first();
            expect(nextButton.text()).toBe('Next');
            expect(nextButton.prop('disabled')).toBeTruthy();
          });
        });
        describe('when on the last page', () => {
          beforeEach(() => {
            const pageIsValid = false;

            wrapper = shallow(
              <WizardPage {...minProps} pageKey="3" pageIsValid={pageIsValid} match={{}}>
                <div>This is page 1</div>
              </WizardPage>,
            );
            buttons = wrapper.find('button');
          });
          it('the complete button is second to last and is disabled', () => {
            const nextButton = buttons.at(1);
            expect(nextButton.text()).toBe('Complete');
            expect(nextButton.prop('disabled')).toBeTruthy();
          });
        });
      });
      describe('when pageIsValid is true', () => {
        beforeEach(() => {
          const continueToNextPage = true;

          wrapper = shallow(
            <WizardPage {...minProps} pageIsValid={continueToNextPage}>
              <div>This is page 1</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('the next button is enabled', () => {
          const nextButton = buttons.at(0);
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
              {...minProps}
              pageKey="2"
              push={mockPush}
              match={{}}
              hasSucceeded={false}
              error={{ message: 'Something bad happened' }}
            >
              <div>This is page 2</div>
            </WizardPage>,
          );
          buttons = wrapper.find('button');
        });
        it('it shows an error alert before its child', () => {
          const childContainer = wrapper.find('div.error-message');
          expect(childContainer.exists()).toBe(true);
          expect(childContainer.first().text()).toBe('<Alert />');
        });
        it('it renders button for cancel, back, next', () => {
          expect(buttons.length).toBe(3);
        });
        it('the cancel button is last and is enabled', () => {
          const cancelButton = buttons.last();
          expect(cancelButton.text()).toBe('Cancel');
          expect(cancelButton.prop('disabled')).toBe(false);
        });
        it('the back button is first and is enabled', () => {
          const backButton = buttons.first();
          expect(backButton.text()).toBe('Back');
          expect(backButton.prop('disabled')).toBe(false);
        });
        it('the next button is second and is enabled', () => {
          const nextButton = buttons.at(1);
          expect(nextButton.text()).toBe('Next');
          expect(nextButton.prop('disabled')).toBe(false);
        });
      });
    });
    describe('when page is not dirty', () => {
      beforeEach(() => {
        mockPush.mockClear();
        wrapper = shallow(
          <WizardPage {...minProps} pageKey="2" push={mockPush} match={{}} hasSucceeded={false} dirty={false}>
            <div>This is page 2</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      it('the previous button is first and is enabled', () => {
        const prevButton = buttons.first();
        expect(prevButton.text()).toBe('Back');
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
        it('submit is not called', () => {
          expect(submit.mock.calls.length).toBe(0);
        });
      });
      it('the cancel button is last and is enabled', () => {
        const prevButton = buttons.last();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      it('the next button is second and is enabled', () => {
        const nextButton = buttons.at(1);
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
          <WizardPage {...minProps} push={mockPush} match={{}}>
            <div>This is page 1</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => mockPush.mockClear());
      it('it starts on the first page', () => {
        expect(wrapper.children().first().text()).toBe('This is page 1');
      });
      it('it renders button for cancel and next', () => {
        expect(buttons.length).toBe(2);
      });
      it('the cancel button is last and is enabled', () => {
        const prevButton = buttons.last();

        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      it('the next button is first and is enabled', () => {
        const nextButton = buttons.first();
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
          <WizardPage {...minProps} pageKey="2" push={mockPush} match={{}}>
            <div>This is page 2</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      it('it shows its child', () => {
        expect(wrapper.children().first().text()).toBe('This is page 2');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is first and is enabled', () => {
        const prevButton = buttons.first();
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
      it('the cancel button is last and is enabled', () => {
        const prevButton = buttons.last();
        expect(prevButton.text()).toBe('Cancel');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      it('the next button is second and is enabled', () => {
        const nextButton = buttons.at(1);
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
          <WizardPage {...minProps} handleSubmit={submit} pageKey="3" push={mockPush} match={{}}>
            <div>This is page 3</div>
          </WizardPage>,
        );
        buttons = wrapper.find('button');
      });
      afterEach(() => {
        submit.mockClear();
      });

      it('it shows its child', () => {
        expect(wrapper.children().first().text()).toBe('This is page 3');
      });
      it('it renders button for cancel, back, next', () => {
        expect(buttons.length).toBe(3);
      });
      it('the back button is first and is enabled', () => {
        const prevButton = buttons.first();
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
      it('the cancel button is last and is enabled', () => {
        const saveButton = buttons.last();
        expect(saveButton.text()).toBe('Cancel');
        expect(saveButton.prop('disabled')).toBe(false);
      });
      it('the Complete button is second and is enabled', () => {
        const nextButton = buttons.at(1);
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
          {...minProps}
          pageList={['page1', 'anotherPage/:foo/:bar']}
          pageKey="page1"
          match={{ params: { foo: 'dvorak' } }}
          push={mockPush}
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
  describe('when there is a canMoveNext prop', () => {
    describe('when canMoveNext is true', () => {
      wrapper = shallow(<WizardPage {...minProps} />);
      const nextButton = wrapper.find('button').last();
      expect(nextButton.text()).toBe('Next');
      expect(nextButton.prop('disabled')).toEqual(false);
    });
    describe('when canMoveNext is false', () => {
      wrapper = shallow(<WizardPage {...minProps} canMoveNext={false} />);
      const nextButton = wrapper.find('button').at(1);
      expect(nextButton.text()).toBe('Next');
      expect(nextButton.prop('disabled')).toEqual(true);
    });
  });
});
