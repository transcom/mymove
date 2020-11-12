/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import WizardNavigation from './WizardNavigation';

describe('WizardNavigation', () => {
  describe('with default props', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };

    const wrapper = mount(<WizardNavigation {...mockProps} />);
    const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
    const backButton = wrapper.find('button[data-testid="wizardBackButton"]');

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
      expect(backButton.length).toBe(1);
      expect(backButton.text()).toEqual('Back');
      expect(nextButton.length).toBe(1);
      expect(nextButton.text()).toEqual('Next');
    });

    it('does not render complete or cancel buttons', () => {
      expect(wrapper.find('button[data-testid="wizardCancelButton"]').length).toBe(0);
      expect(wrapper.find('button[data-testid="wizardCompleteButton"]').length).toBe(0);
    });

    it('hooks up the onClick handlers', () => {
      nextButton.simulate('click');
      expect(mockProps.onNextClick).toHaveBeenCalled();
      backButton.simulate('click');
      expect(mockProps.onBackClick).toHaveBeenCalled();
    });
  });

  describe('if the next button is disabled', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };

    const wrapper = mount(<WizardNavigation {...mockProps} disableNext />);
    const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
    it('the onNextClick handler is not called', () => {
      nextButton.simulate('click');
      expect(mockProps.onNextClick).not.toHaveBeenCalled();
    });
  });

  describe('on the first page', () => {
    const wrapper = mount(<WizardNavigation isFirstPage />);

    it('doesn’t show the back button', () => {
      expect(wrapper.find('[data-testid="wizardBackButton"]').length).toBe(0);
    });
  });

  describe('on the last page', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };
    const wrapper = mount(<WizardNavigation {...mockProps} isLastPage />);
    const nextButton = wrapper.find('button[data-testid="wizardNextButton"]');
    const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');

    it('shows the complete button', () => {
      expect(nextButton.length).toBe(0);
      expect(completeButton.length).toBe(1);
      expect(completeButton.text()).toEqual('Complete');
    });

    it('hooks up the onClick handlers', () => {
      completeButton.simulate('click');
      expect(mockProps.onNextClick).toHaveBeenCalled();
    });
  });

  describe('if the complete button is disabled on the last page', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };

    const wrapper = mount(<WizardNavigation {...mockProps} disableNext isLastPage />);
    const completeButton = wrapper.find('button[data-testid="wizardCompleteButton"]');

    it('the onNextClick handler is not called', () => {
      completeButton.simulate('click');
      expect(mockProps.onNextClick).not.toHaveBeenCalled();
    });
  });

  describe('if Finish Later is an option', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };
    const wrapper = mount(<WizardNavigation showFinishLater {...mockProps} />);
    const finishLaterButton = wrapper.find('button[data-testid="wizardCancelButton"]');

    it('shows the finish later button', () => {
      expect(finishLaterButton.length).toBe(1);
    });

    it('hooks up the onClick handlers', () => {
      finishLaterButton.simulate('click');
      expect(mockProps.onCancelClick).toHaveBeenCalled();
    });
  });

  describe('if in edit mode', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };
    const wrapper = mount(<WizardNavigation editMode {...mockProps} />);
    const saveButton = wrapper.find('button[data-testid="wizardNextButton"]');
    const cancelButton = wrapper.find('button[data-testid="wizardCancelButton"]');

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
      expect(saveButton.length).toBe(1);
      expect(saveButton.text()).toEqual('Save');
      expect(cancelButton.length).toBe(1);
      expect(cancelButton.text()).toEqual('Cancel');
    });

    it('does not render complete or back buttons', () => {
      expect(wrapper.find('button[data-testid="wizardBackButton"]').length).toBe(0);
      expect(wrapper.find('button[data-testid="wizardCompleteButton"]').length).toBe(0);
    });

    it('hooks up the onClick handlers', () => {
      saveButton.simulate('click');
      expect(mockProps.onNextClick).toHaveBeenCalled();
      cancelButton.simulate('click');
      expect(mockProps.onCancelClick).toHaveBeenCalled();
    });
  });

  describe('if in readOnly mode', () => {
    const mockProps = {
      onBackClick: jest.fn(),
      onNextClick: jest.fn(),
      onCancelClick: jest.fn(),
    };
    const wrapper = mount(<WizardNavigation readOnly {...mockProps} />);
    const cancelButton = wrapper.find('button[data-testid="wizardCancelButton"]');

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
      expect(cancelButton.length).toBe(1);
      expect(cancelButton.text()).toEqual('Return home');
    });

    it('only renders the return home button', () => {
      expect(wrapper.find('button[data-testid="wizardBackButton"]').length).toBe(0);
      expect(wrapper.find('button[data-testid="wizardNextButton"]').length).toBe(0);
      expect(wrapper.find('button[data-testid="wizardCompleteButton"]').length).toBe(0);
    });

    it('hooks up the onClick handlers', () => {
      cancelButton.simulate('click');
      expect(mockProps.onCancelClick).toHaveBeenCalled();
    });
  });
});
