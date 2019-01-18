import { getDetailComponent } from './DetailsHelper';
import { DefaultDetails } from './DefaultDetails';
import { Code105Details } from './Code105Details';

describe('testing getDetailComponent()', () => {
  describe('returns default details component', () => {
    const DetailComponent = getDetailComponent();

    it('renders without crashing', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });

  describe('returns default details component with feature flag off', () => {
    //pass in known code item with feature flag off
    let DetailComponent = getDetailComponent('105', false);
    it('renders default details without crashing', () => {
      expect(DetailComponent).toBe(DefaultDetails);
    });
  });

  describe('returns 105B/E details component with feature flag on', () => {
    let DetailComponent = getDetailComponent('105', true);
    it('renders 105 details without crashing', () => {
      expect(DetailComponent).toBe(Code105Details);
    });
  });
});
