import { shape, string } from 'prop-types';

export const CustomerSupportRemarkShape = shape({
  moveID: string.isRequired,
  officeUserID: string.isRequired,
  content: string.isRequired,
  officeUserFirstName: string.isRequired,
  officeUserLastName: string.isRequired,
  updatedAt: string.isRequired,
  createdAt: string.isRequired,
});

export default {
  CustomerSupportRemarkShape,
};
