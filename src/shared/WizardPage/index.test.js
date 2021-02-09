import React from 'react';
import { mount } from 'enzyme';
import { act } from 'react-dom/test-utils';

import { WizardPage } from 'shared/WizardPage';

describe('the WizardPage component', () => {
  const minProps = {
    handleSubmit: jest.fn(),
    pageList: ['1', '2', '3'],
    pageKey: '1',
    match: {},
  };

  const middleFlowProps = {
    handleSubmit: jest.fn(),
    pageList: ['1', '2', '3'],
    pageKey: '2',
    push: jest.fn(),
    match: {},
  };

  describe('with minimum props', () => {
    const wrapper = mount(<WizardPage {...minProps} />);
    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('renders navigation', () => {
      expect(wrapper.find('WizardNavigation').exists()).toBe(true);
    });
  });

  describe('navigation', () => {
    describe('on the first page', () => {
      const mockPush = jest.fn();
      const wrapper = mount(
        <WizardPage {...minProps} push={mockPush}>
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
        it('push gets the next page', async () => {
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          await act(async () => {
            nextButton.simulate('click');
          });
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('2');
        });
      });
    });

    describe('when on the middle page', () => {
      const mockPush = jest.fn();
      const mockSubmit = jest.fn();
      const wrapper = mount(
        <WizardPage {...middleFlowProps} push={mockPush} handleSubmit={mockSubmit}>
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
          mockPush.mockClear();
          mockSubmit.mockClear();
          const prevButton = wrapper.find('button[data-testid="wizardBackButton"]');
          prevButton.simulate('click');
        });

        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('1');
        });

        it('submit is not called', () => {
          expect(mockSubmit.mock.calls.length).toBe(0);
        });
      });

      describe('when the next button is clicked', () => {
        beforeEach(() => {
          mockPush.mockClear();
          mockSubmit.mockClear();
          const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
          nextButton.simulate('click');
        });

        it('push gets the next page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('3');
        });

        it('submit is called', () => {
          expect(mockSubmit.mock.calls.length).toBe(1);
        });
      });
    });

    describe('on the last page', () => {
      const mockPush = jest.fn();
      const mockSubmit = jest.fn();
      const wrapper = mount(
        <WizardPage {...minProps} pageKey="3" handleSubmit={mockSubmit} push={mockPush}>
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

        it('push gets the prev page', () => {
          expect(mockPush.mock.calls.length).toBe(1);
          expect(mockPush.mock.calls[0][0]).toBe('2');
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
        const wrapper = mount(
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
        const wrapper = mount(
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
        const wrapper = mount(
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
        const wrapper = mount(
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
    const wrapper = mount(
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
    const mockPush = jest.fn();

    const wrapper = mount(
      <WizardPage {...minProps} pageKey="2" push={mockPush} handleSubmit={mockSubmit} dirty={false}>
        <div>This is page 2</div>
      </WizardPage>,
    );

    describe('when the next button is clicked', () => {
      beforeEach(() => {
        mockSubmit.mockClear();
        mockPush.mockClear();
        const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
        nextButton.simulate('click');
      });

      it('push gets the next page', () => {
        expect(mockPush.mock.calls.length).toBe(1);
        expect(mockPush.mock.calls[0][0]).toBe('3');
      });

      it('submit is not called', () => {
        expect(mockSubmit.mock.calls.length).toBe(0);
      });
    });
  });

  describe('when there is an additionalParams prop', () => {
    const mockPush = jest.fn();
    const wrapper = mount(
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

    it('clicking the next button calls the handler with additionalParams', async () => {
      const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
      await act(async () => {
        nextButton.simulate('click');
      });

      expect(mockPush).toHaveBeenCalledWith('anotherPage/dvorak/querty');
    });
  });

  describe('when canMoveNext is true', () => {
    const wrapper = mount(<WizardPage {...minProps} />);
    it('the WizardNavigation disableNext prop is false', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('disableNext')).toEqual(false);
    });
  });

  describe('when canMoveNext is false', () => {
    const wrapper = mount(<WizardPage {...minProps} canMoveNext={false} />);
    it('the WizardNavigation disableNext prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('disableNext')).toEqual(true);
    });
  });

  describe('when hideBackBtn is true on middle page', () => {
    const wrapper = mount(<WizardPage {...middleFlowProps} hideBackBtn />);

    it('the WizardNavigation isFirstPage prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('isFirstPage')).toEqual(true);
    });
  });

  describe('when showFinishLaterBtn is true', () => {
    const wrapper = mount(<WizardPage {...middleFlowProps} showFinishLaterBtn />);

    it('the WizardNavigation showFinishLater prop is true', () => {
      const navigation = wrapper.find('WizardNavigation');
      expect(navigation.prop('showFinishLater')).toEqual(true);
    });
  });
});
