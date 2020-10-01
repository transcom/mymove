import React from 'react';
import { shallow } from 'enzyme';

import { WizardPage } from 'shared/WizardPage';

describe('given a WizardPage', () => {
  let wrapper;
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
          });

          it('it renders buttons for next', () => {
            const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
            expect(nextButton.text()).toBe('Next');
          });

          it('does not render a back button', () => {
            const backButton = wrapper.find('[data-testid="wizardBackButton"]');
            expect(backButton.exists()).toBe(false);
          });

          it('the next button is disabled', () => {
            const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
          });

          it('it renders buttons for back and complete', () => {
            const backButton = wrapper.find('[data-testid="wizardBackButton"]');
            expect(backButton.text()).toBe('Back');
            const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
            expect(completeButton.text()).toBe('Complete');
          });

          it('does not render a next button', () => {
            const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
            expect(nextButton.exists()).toBe(false);
          });

          it('the complete button is disabled', () => {
            const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
            expect(completeButton.prop('disabled')).toBeTruthy();
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
        });

        it('the next button is enabled', () => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
        });

        it('it shows an error alert before its child', () => {
          const childContainer = wrapper.find('div.error-message');
          expect(childContainer.exists()).toBe(true);
          expect(childContainer.first().text()).toBe('<Alert />');
        });

        it('it renders button for back, next', () => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
          expect(nextButton.text()).toBe('Next');
          const backButton = wrapper.find('[data-testid="wizardBackButton"]');
          expect(backButton.text()).toBe('Back');
        });

        it('does not render a complete button', () => {
          const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
          expect(completeButton.exists()).toBe(false);
        });

        it('the back button is enabled', () => {
          const backButton = wrapper.find('[data-testid="wizardBackButton"]');
          expect(backButton.prop('disabled')).toBe(false);
        });

        it('the next button is enabled', () => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
      });
      it('the previous button is enabled', () => {
        const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(prevButton.prop('disabled')).toBe(false);
      });

      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
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
      it('the next button is enabled', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.prop('disabled')).toBeFalsy();
      });
      describe('when the next button is clicked', () => {
        beforeEach(() => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
      });
      afterEach(() => mockPush.mockClear());
      it('it starts on the first page', () => {
        expect(wrapper.children().first().text()).toBe('This is page 1');
      });
      it('it renders button for next', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.text()).toBe('Next');
      });
      it('the next button is enabled', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.prop('disabled')).toBeFalsy();
      });

      describe('when the next button is clicked', () => {
        beforeEach(() => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
      });
      it('it shows its child', () => {
        expect(wrapper.children().first().text()).toBe('This is page 2');
      });
      it('it renders button for back, next', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.text()).toBe('Next');
        const backButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(backButton.text()).toBe('Back');
      });
      it('the back button is enabled', () => {
        const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(prevButton.prop('disabled')).toBe(false);
      });

      describe('when the back button is clicked', () => {
        beforeEach(() => {
          const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
          prevButton.simulate('click');
        });
        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('1');
        });
      });
      it('the next button is enabled', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.prop('disabled')).toBeFalsy();
      });
      describe('when the next button is clicked', () => {
        beforeEach(() => {
          const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
      });
      afterEach(() => {
        submit.mockClear();
      });

      it('it shows its child', () => {
        expect(wrapper.children().first().text()).toBe('This is page 3');
      });
      it('it renders button for back and complete', () => {
        const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
        expect(completeButton.text()).toBe('Complete');
        const backButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(backButton.text()).toBe('Back');
      });

      it('the back button is enabled', () => {
        const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(prevButton.prop('disabled')).toBe(false);
      });
      describe('when the back button is clicked', () => {
        beforeEach(() => {
          const prevButton = wrapper.find('[data-testid="wizardBackButton"]');
          prevButton.simulate('click');
        });
        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('2');
        });
      });
      it('the Complete button is enabled', () => {
        const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
        expect(completeButton.prop('disabled')).toBeFalsy();
      });
      describe('when the complete button is clicked', () => {
        beforeEach(() => {
          const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
          completeButton.simulate('click');
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
    });
    describe('when the next button is clicked', () => {
      beforeEach(() => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
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
      const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
      expect(nextButton.prop('disabled')).toEqual(false);
    });
    describe('when canMoveNext is false', () => {
      wrapper = shallow(<WizardPage {...minProps} canMoveNext={false} />);
      const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
      expect(nextButton.prop('disabled')).toEqual(true);
    });
  });
});
