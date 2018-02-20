import React from 'react';
import { mount, shallow } from 'enzyme';
import Wizard from 'shared/Wizard';

describe('Given a wizard with 3 pages', () => {
  let wrapper;

  beforeEach(() => {
    wrapper = shallow(
      <Wizard>
        <div>This is page 1</div>
        <div>This is page 2</div>
        <div>This is page 3</div>
      </Wizard>,
    );
  });
  it('it starts on the first page', () => {
    const divs = wrapper.find('div');
    expect(divs.length).toBe(1);
    expect(divs.first().text()).toBe('This is page 1');
  });
  describe('When on the first page', () => {
    it('it only renders a next button', () => {
      const buttons = wrapper.find('button');
      expect(buttons.length).toBe(1);
      expect(buttons.first().text()).toBe('Next');
    });
  });
  describe('when the next button is clicked', () => {
    beforeEach(() => {
      const firstButton = wrapper.find('button').first();
      firstButton.simulate('click');
    });
    it('it shows the second page', () => {
      const divs = wrapper.find('div');
      expect(divs.length).toBe(1);
      expect(divs.first().text()).toBe('This is page 2');
    });
    it('it shows the prev and next buttons', () => {
      const buttons = wrapper.find('button');
      expect(buttons.length).toBe(2);
      expect(buttons.first().text()).toBe('Prev');
      expect(buttons.first().hasClass('usa-button-secondary')).toBe(true);
      expect(buttons.last().text()).toBe('Next');
      expect(buttons.last().hasClass('usa-button-secondary')).toBe(false);
    });
    describe('when the next button is clicked a second time', () => {
      beforeEach(() => {
        wrapper
          .find('button')
          .last()
          .simulate('click');
      });
      it('it shows the second page', () => {
        const divs = wrapper.find('div');
        expect(divs.length).toBe(1);
        expect(divs.first().text()).toBe('This is page 3');
      });
      it('it shows only the prev buttons', () => {
        const buttons = wrapper.find('button');
        expect(buttons.length).toBe(1);
        expect(buttons.first().text()).toBe('Prev');
        expect(buttons.first().hasClass('usa-button-secondary')).toBe(false);
      });
    });
    describe('when the prev button is clicked', () => {
      beforeEach(() => {
        wrapper
          .find('button')
          .first()
          .simulate('click');
      });
      it('it shows the first page', () => {
        const divs = wrapper.find('div');
        expect(divs.length).toBe(1);
        expect(divs.first().text()).toBe('This is page 1');
      });
      it('it shows only the next buttons', () => {
        const buttons = wrapper.find('button');
        expect(buttons.length).toBe(1);
        expect(buttons.first().text()).toBe('Next');
        expect(buttons.first().hasClass('usa-button-secondary')).toBe(false);
      });
    });
  });
});
