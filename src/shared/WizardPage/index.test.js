import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';

import { WizardPage } from 'shared/WizardPage';
import { MockRouting } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(() => {
  jest.resetAllMocks();
});

const mountWithRouting = (ui) => {
  return mount(<MockRouting>{ui}</MockRouting>);
};

describe('the WizardPage component', () => {
  const minProps = {
    handleSubmit: jest.fn(),
    pageList: ['1', '2', '3'],
    pageKey: '1',
  };

  const middleFlowProps = {
    handleSubmit: jest.fn(),
    pageList: ['1', '2', '3'],
    pageKey: '2',
  };

  describe('with minimum props', () => {
    const wrapper = mountWithRouting(<WizardPage {...minProps} />);
    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('renders navigation', () => {
      expect(wrapper.find('WizardNavigation').exists()).toBe(true);
    });
  });

  describe('navigation', () => {
    describe('on the first page', () => {
      const wrapper = mountWithRouting(
        <WizardPage {...minProps}>
          <div>This is page 1</div>
        </WizardPage>,
      );

      it('it renders buttons for next', () => {
        const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
        expect(nextButton.text()).toBe('Next');
        expect(nextButton.prop('disabled')).toBeFalsy();
      });

      it('does not render a back button', () => {
        const backButton = wrapper.find('[data-testid="wizardBackButton"]');
        expect(backButton.exists()).toBe(false);
      });

      describe('when the next button is clicked', () => {
        it('navigate gets the next page', async () => {
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          await act(async () => {
            nextButton.simulate('click');
          });
          expect(mockNavigate.mock.calls.length).toBe(1);
          expect(mockNavigate.mock.calls[0][0]).toBe('2');
        });
      });
    });

    describe('when on the middle page', () => {
      const mockSubmit = jest.fn();
      const wrapper = mountWithRouting(
        <WizardPage {...middleFlowProps} handleSubmit={mockSubmit}>
          <div>This is page 2</div>
        </WizardPage>,
      );

      it('it renders button for back and next', () => {
        const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
        expect(nextButton.text()).toBe('Next');
        const backButton = wrapper.find('button[data-testid="wizardBackButton"]');
        expect(backButton.text()).toBe('Back');
      });

      it('does not render a complete button', () => {
        const completeButton = wrapper.find('[data-testid="wizardCompleteButton"]');
        expect(completeButton.exists()).toBe(false);
      });

      it('the back button is enabled', () => {
        const backButton = wrapper.find('button[data-testid="wizardBackButton"]');
        expect(backButton.prop('disabled')).toBeFalsy();
      });

      describe('when the prev button is clicked', () => {
        beforeEach(() => {
          const prevButton = wrapper.find('button[data-testid="wizardBackButton"]');
          prevButton.simulate('click');
        });

        it('navigate gets the prev page', () => {
          expect(mockNavigate.mock.calls.length).toBe(1);
          expect(mockNavigate.mock.calls[0][0]).toBe('1');
        });

        it('submit is not called', () => {
          expect(mockSubmit.mock.calls.length).toBe(0);
        });
      });

      describe('when the next button is clicked', () => {
        beforeEach(() => {
          mockNavigate.mockClear();
          mockSubmit.mockClear();
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          nextButton.simulate('click');
        });

        it('navigate gets the next page', () => {
          expect(mockNavigate.mock.calls.length).toBe(1);
          expect(mockNavigate.mock.calls[0][0]).toBe('3');
        });

        it('submit is called', () => {
          expect(mockSubmit.mock.calls.length).toBe(1);
        });
      });
    });

    describe('on the last page', () => {
      const mockSubmit = jest.fn();
      const wrapper = mountWithRouting(
        <WizardPage {...minProps} pageKey="3" handleSubmit={mockSubmit}>
          <div>This is page 3</div>
        </WizardPage>,
      );

      it('it renders buttons for back and complete', () => {
        const backButton = wrapper.find('button[data-testid="wizardBackButton"]');
        expect(backButton.text()).toBe('Back');
        expect(backButton.prop('disabled')).toBeFalsy();

        const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');
        expect(completeButton.text()).toBe('Complete');
      });

      it('does not render a next button', () => {
        const nextButton = wrapper.find('[data-testid="wizardNextButton"]');
        expect(nextButton.exists()).toBe(false);
      });

      describe('when the back button is clicked', () => {
        beforeEach(() => {
          const prevButton = wrapper.find('button[data-testid="wizardBackButton"]');
          prevButton.simulate('click');
        });

        it('navigate gets the prev page', () => {
          expect(mockNavigate.mock.calls.length).toBe(1);
          expect(mockNavigate.mock.calls[0][0]).toBe('2');
        });
      });

      describe('when the complete button is clicked', () => {
        beforeEach(() => {
          const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');
          completeButton.simulate('click');
        });

        it('submit is called', () => {
          expect(mockSubmit.mock.calls.length).toBe(1);
        });
      });
    });
  });

  describe('when the pageIsValid prop', () => {
    describe('is false', () => {
      describe('on the first page', () => {
        const wrapper = mountWithRouting(
          <WizardPage {...minProps} pageIsValid={false}>
            <div>This is page 1</div>
          </WizardPage>,
        );

        it('the next button is disabled', () => {
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          expect(nextButton.prop('disabled')).toBeTruthy();
        });
      });

      describe('on the last page', () => {
        const wrapper = mountWithRouting(
          <WizardPage {...minProps} pageKey="3" pageIsValid={false}>
            <div>This is page 3</div>
          </WizardPage>,
        );

        it('the complete button is disabled', () => {
          const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');
          expect(completeButton.prop('disabled')).toBeTruthy();
        });
      });
    });

    describe('is true', () => {
      describe('on the first page', () => {
        const wrapper = mountWithRouting(
          <WizardPage {...minProps} pageIsValid={true}>
            <div>This is page 1</div>
          </WizardPage>,
        );

        it('the next button is enabled', () => {
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          expect(nextButton.prop('disabled')).toBeFalsy();
        });
      });

      describe('on the last page', () => {
        const wrapper = mountWithRouting(
          <WizardPage {...minProps} pageKey="3" pageIsValid={true}>
            <div>This is page 3</div>
          </WizardPage>,
        );

        it('the complete button is enabled', () => {
          const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');
          expect(completeButton.prop('disabled')).toBeFalsy();
        });
      });
    });
  });

  describe('when there is an error', () => {
    const wrapper = mountWithRouting(
      <WizardPage {...middleFlowProps} error={{ message: 'Something bad happened' }}>
        <div>This is page 2</div>
      </WizardPage>,
    );

    it('it renders an error alert', () => {
      expect(wrapper.find('Alert[type="error"]').exists()).toBe(true);
    });
  });

  describe('when page is not dirty', () => {
    const mockSubmit = jest.fn();

    const wrapper = mountWithRouting(
      <WizardPage {...minProps} pageKey="2" handleSubmit={mockSubmit} dirty={false}>
        <div>This is page 2</div>
      </WizardPage>,
    );

    describe('when the next button is clicked', () => {
      beforeEach(() => {
        const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
        nextButton.simulate('click');
      });

      it('navigate gets the next page', () => {
        expect(mockNavigate.mock.calls.length).toBe(1);
        expect(mockNavigate.mock.calls[0][0]).toBe('3');
      });

      it('submit is not called', () => {
        expect(mockSubmit.mock.calls.length).toBe(0);
      });
    });
  });

  describe('when there is an additionalParams prop', () => {
    const wrapper = mountWithRouting(
      <WizardPage
        {...minProps}
        pageList={['page1', 'anotherPage/:foo/:bar']}
        pageKey="page1"
        match={{ params: { foo: 'dvorak' } }}
        additionalParams={{ bar: 'querty' }}
      >
        <div>This is page 1</div>
      </WizardPage>,
    );

    it('clicking the next button calls the handler with additionalParams', async () => {
      const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
      await act(async () => {
        nextButton.simulate('click');
      });

      expect(mockNavigate).toHaveBeenCalledWith('anotherPage/dvorak/querty');
    });
  });

  describe('when canMoveNext is true', () => {
    const wrapper = mountWithRouting(<WizardPage {...minProps} />);
    it('the WizardNavigation disableNext prop is false', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('disableNext')).toEqual(false);
    });
  });

  describe('when canMoveNext is false', () => {
    const wrapper = mountWithRouting(<WizardPage {...minProps} canMoveNext={false} />);
    it('the WizardNavigation disableNext prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('disableNext')).toEqual(true);
    });
  });

  describe('when hideBackBtn is true on middle page', () => {
    const wrapper = mountWithRouting(<WizardPage {...middleFlowProps} hideBackBtn />);

    it('the WizardNavigation isFirstPage prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('isFirstPage')).toEqual(true);
    });
  });

  describe('when showFinishLaterBtn is true', () => {
    const wrapper = mountWithRouting(<WizardPage {...middleFlowProps} showFinishLaterBtn />);

    it('the WizardNavigation showFinishLater prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('showFinishLater')).toEqual(true);
    });
  });
});
