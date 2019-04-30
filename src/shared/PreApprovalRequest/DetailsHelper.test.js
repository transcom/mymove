import { getFormComponent, isRobustAccessorial } from './DetailsHelper';
import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';
import { Code35Form } from './Code35Form';
import { Code226Form } from './Code226Form';

let featureFlag = false;
let initialValuesWithoutCrateDimensions = {};
let initialValuesWithCrateDimensions = { crate_dimensions: true };
describe('testing getFormComponent()', () => {
  describe('returns default form component', () => {
    const FormComponent = getFormComponent();

    it('for undefined values', () => {
      expect(FormComponent).toBe(DefaultForm);
    });
  });

  describe('returns default form component with feature flag off', () => {
    //pass in known code item with feature flag off
    featureFlag = false;

    let FormComponent;
    it('for code 105', () => {
      FormComponent = getFormComponent('105', featureFlag, initialValuesWithCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });

    //testing for non-existing code
    it('for code 4A', () => {
      FormComponent = getFormComponent('4A', featureFlag, initialValuesWithCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });

    it('for code 105D', () => {
      FormComponent = getFormComponent('105D', featureFlag, initialValuesWithCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });

    it('for code 105D', () => {
      FormComponent = getFormComponent('105D', featureFlag, null);
      expect(FormComponent).toBe(DefaultForm);
    });
  });

  describe('returns 105B/E form component with feature flag on', () => {
    featureFlag = true;
    let FormComponent;

    it('for code 105B', () => {
      FormComponent = getFormComponent('105B', featureFlag, initialValuesWithCrateDimensions);
      expect(FormComponent).toBe(Code105Form);
    });

    it('for code 105E', () => {
      FormComponent = getFormComponent('105E', featureFlag, initialValuesWithCrateDimensions);
      expect(FormComponent).toBe(Code105Form);
    });
  });

  describe('returns 35A form component with feature flag on', () => {
    featureFlag = true;

    it('for code 35A', () => {
      let FormComponent = getFormComponent('35A', featureFlag, { estimate_amount_cents: true });
      expect(FormComponent).toBe(Code35Form);
    });
  });

  describe('returns 226A form component with feature flag on', () => {
    featureFlag = true;

    it('for code 226A', () => {
      let FormComponent = getFormComponent('226A', featureFlag, { actual_amount_cents: true });
      expect(FormComponent).toBe(Code226Form);
    });
  });

  describe('returns default form component with feature flag on', () => {
    featureFlag = true;
    let FormComponent;

    it('for code 105D', () => {
      FormComponent = getFormComponent('105D', featureFlag);
      expect(FormComponent).toBe(DefaultForm);
    });

    it('for code 105B without crate dimensions', () => {
      FormComponent = getFormComponent('105B', featureFlag, initialValuesWithoutCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });

    it('for code 105E without crate dimensions', () => {
      FormComponent = getFormComponent('105E', featureFlag, initialValuesWithoutCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });

    it('for code 226A without crate dimensions', () => {
      FormComponent = getFormComponent('226A', featureFlag, initialValuesWithoutCrateDimensions);
      expect(FormComponent).toBe(DefaultForm);
    });
  });
});

describe('preApprovals', () => {
  describe('isNewAccessorial 105B or 105E', () => {
    it('should return true if new accessorial, false if old accessorial', () => {
      const item105BOld = { tariff400ng_item: { code: '105B' } };
      const item105BNew = { tariff400ng_item: { code: '105B' }, crate_dimensions: { length: 1, height: 1, width: 1 } };
      const item105EOld = { tariff400ng_item: { code: '105E' } };
      const item105ENew = { tariff400ng_item: { code: '105E' }, crate_dimensions: { length: 1, height: 1, width: 1 } };

      const itemNull = null;

      expect(isRobustAccessorial(item105BOld)).toEqual(false);
      expect(isRobustAccessorial(item105BNew)).toEqual(true);
      expect(isRobustAccessorial(item105EOld)).toEqual(false);
      expect(isRobustAccessorial(item105ENew)).toEqual(true);
      expect(isRobustAccessorial(itemNull)).toEqual(false);
    });
  });

  describe('isNewAccessorial 35A', () => {
    it('should return true if new accessorial, false if old accessorial', () => {
      const item35AOld = { tariff400ng_item: { code: '35A' } };
      const item35ANew = { tariff400ng_item: { code: '35A' }, estimate_amount_cents: 1 };
      const itemNull = null;

      expect(isRobustAccessorial(item35AOld)).toEqual(false);
      expect(isRobustAccessorial(item35ANew)).toEqual(true);
      expect(isRobustAccessorial(itemNull)).toEqual(false);
    });
  });

  describe('isNewAccessorial 226A', () => {
    it('should return true if new accessorial, false if old accessorial', () => {
      const item226AOld = { tariff400ng_item: { code: '226A' } };
      const item226ANew = { tariff400ng_item: { code: '226A' }, actual_amount_cents: 1 };
      const itemNull = null;

      expect(isRobustAccessorial(item226AOld)).toEqual(false);
      expect(isRobustAccessorial(item226ANew)).toEqual(true);
      expect(isRobustAccessorial(itemNull)).toEqual(false);
    });
  });

  describe('isNewAccessorial 125', () => {
    it('should return true if new accessorial, false if old accessorial', () => {
      const item125Old = { tariff400ng_item: { code: '125A' } };
      const item125New = { tariff400ng_item: { code: '125A' }, address: {} };
      const itemNull = null;

      expect(isRobustAccessorial(item125Old)).toEqual(false);
      expect(isRobustAccessorial(item125New)).toEqual(true);
      expect(isRobustAccessorial(itemNull)).toEqual(false);
    });
  });
});
