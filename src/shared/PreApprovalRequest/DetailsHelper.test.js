import { getFormComponent } from './DetailsHelper';
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

    let FormComponent = getFormComponent('105', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105B', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105B', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105E', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105E', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    //testing for non-existing code
    FormComponent = getFormComponent('4A', featureFlag, initialValuesWithCrateDimensions);
    it('for code 4A', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105D', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105D', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105D', featureFlag, null);
    it('for code 105D', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('35A', featureFlag);
    it('for code 35A', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('226A', featureFlag);
    it('for code 226A', () => {
      expect(FormComponent).toBe(DefaultForm);
    });
  });

  describe('returns 105B/E form component with feature flag on', () => {
    featureFlag = true;

    let FormComponent = getFormComponent('105', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105', () => {
      expect(FormComponent).toBe(Code105Form);
    });

    FormComponent = getFormComponent('105B', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105B', () => {
      expect(FormComponent).toBe(Code105Form);
    });

    FormComponent = getFormComponent('105E', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105E', () => {
      expect(FormComponent).toBe(Code105Form);
    });
  });

  describe('returns 35A form component with feature flag on', () => {
    featureFlag = true;

    let FormComponent = getFormComponent('35A', featureFlag, { estimate_amount_cents: true });
    it('for code 35A', () => {
      expect(FormComponent).toBe(Code35Form);
    });
  });

  describe('returns 226A form component with feature flag on', () => {
    featureFlag = true;

    let FormComponent = getFormComponent('226A', featureFlag, initialValuesWithCrateDimensions);
    it('for code 226A', () => {
      expect(FormComponent).toBe(Code226Form);
    });
  });

  describe('returns default form component with feature flag on', () => {
    featureFlag = true;

    let FormComponent = getFormComponent('105D', featureFlag);
    it('for code 105D', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105B', featureFlag, initialValuesWithoutCrateDimensions);
    it('for code 105B without crate dimensions', () => {
      expect(FormComponent).toBe(DefaultForm);
    });

    FormComponent = getFormComponent('105E', featureFlag, initialValuesWithoutCrateDimensions);
    it('for code 105E without crate dimensions', () => {
      expect(FormComponent).toBe(DefaultForm);
    });
  });
});
