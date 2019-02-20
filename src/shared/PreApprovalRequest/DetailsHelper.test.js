import { getFormComponent } from './DetailsHelper';
import { DefaultForm } from './DefaultForm';
import { Code105Form } from './Code105Form';

let featureFlag = false;
describe('testing getFormComponent()', () => {
let initialValuesWithoutCrateDimensions = {};
let initialValuesWithCrateDimensions = { crate_dimensions: true };
describe('testing getDetailComponent()', () => {
  describe('returns default details component', () => {
    const DetailComponent = getFormComponent();

    it('for undefined values', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });
  });

  describe('returns default details component with feature flag off', () => {
    //pass in known code item with feature flag off
    featureFlag = false;

    let DetailComponent = getFormComponent('105', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });

    DetailComponent = getFormComponent('105B', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105B', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });

    DetailComponent = getFormComponent('105E', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105E', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });

    //testing for non-existing code
    DetailComponent = getFormComponent('4A', featureFlag, initialValuesWithCrateDimensions);
    it('for code 4A', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });

    DetailComponent = getFormComponent('105D', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });

    DetailComponent = getFormComponent('105D', featureFlag, null);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultForm);
    });
  });

  describe('returns 105B/E details component with feature flag on', () => {
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

  describe('returns default details component with feature flag on', () => {
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
