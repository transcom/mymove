import { getDetailComponent } from './DetailsHelper';
import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';

let featureFlag = false;
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

    let DetailComponent = getDetailComponent('105', featureFlag);
    it('for code 105', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105B', featureFlag);
    it('for code 105B', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105E', featureFlag);
    it('for code 105E', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    //testing for non-existing code
    DetailComponent = getDetailComponent('4A', featureFlag);
    it('for code 4A', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });

    DetailComponent = getDetailComponent('105D', featureFlag);
    it('for code 105D', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });

  describe('returns 105B/E details component with feature flag on', () => {
    featureFlag = true;

    let DetailComponent = getDetailComponent('105', featureFlag);
    it('for code 105', () => {
      expect(DetailComponent).toBe(Code105Details);
    });

    DetailComponent = getDetailComponent('105B', featureFlag);
    it('for code 105B', () => {
      expect(DetailComponent).toBe(Code105Details);
    });

    DetailComponent = getDetailComponent('105E', featureFlag);
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
  });
});
