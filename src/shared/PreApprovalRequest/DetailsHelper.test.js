import { getDetailComponent } from './DetailsHelper';
import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';

let featureFlag = false;
let initialValuesWithoutCrateDimensions = {};
let initialValuesWithCrateDimensions = { crate_dimensions: true };
describe('testing getDetailComponent()', () => {
  describe('returns default details component', () => {
    const DetailComponent = getDetailComponent();

    it('for undefined values', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });

  describe('returns default details component with feature flag off', () => {
    //pass in known code item with feature flag off
    featureFlag = false;

    let DetailComponent = getDetailComponent('105', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105B', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105B', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105E', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105E', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    //testing for non-existing code
    DetailComponent = getDetailComponent('4A', featureFlag, initialValuesWithCrateDimensions);
    it('for code 4A', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105D', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105D', featureFlag, null);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });

  describe('returns 105B/E details component with feature flag on', () => {
    featureFlag = true;

    let DetailComponent = getDetailComponent('105', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105', () => {
      expect(DetailComponent).toBe(Code105Details);
    });

    DetailComponent = getDetailComponent('105B', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105B', () => {
      expect(DetailComponent).toBe(Code105Details);
    });

    DetailComponent = getDetailComponent('105E', featureFlag, initialValuesWithCrateDimensions);
    it('for code 105E', () => {
      expect(DetailComponent).toBe(Code105Details);
    });

    DetailComponent = getDetailComponent('105E', featureFlag, null);
    it('for code 105E', () => {
      expect(DetailComponent).toBe(Code105Details);
    });
  });

  describe('returns default details component with feature flag on', () => {
    featureFlag = true;

    let DetailComponent = getDetailComponent('105D', featureFlag);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105B', featureFlag, initialValuesWithoutCrateDimensions);
    it('for code 105B without crate dimensions', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105E', featureFlag, initialValuesWithoutCrateDimensions);
    it('for code 105E without crate dimensions', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105E', featureFlag, null);
    it('for code 105E without crate dimensions', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });
});
