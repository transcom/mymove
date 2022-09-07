import { shape, string, number, bool } from 'prop-types';

export const PWSViolationShape = shape({
  id: string.isRequired,
  displayOrder: number,
  paragraphNumber: string,
  title: string.isRequired,
  category: string.isRequired,
  subCategory: string.isRequired,
  requirementSummary: string,
  requirementStatement: string,
  isKpi: bool,
  additionalDataElem: string,
});

export default {
  PWSViolationShape,
};
