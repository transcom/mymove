import * as Yup from 'yup';
import PropTypes from 'prop-types';

export const uploadShape = Yup.object({
  id: Yup.string(),
  createdAt: Yup.string(),
  bytes: Yup.number(),
  url: Yup.string(),
  filename: Yup.string(),
  contentType: Yup.string(),
});

export const ExistingUploadsShape = PropTypes.arrayOf(
  PropTypes.shape({
    id: PropTypes.string.isRequired,
    createdAt: PropTypes.string.isRequired,
    bytes: PropTypes.number.isRequired,
    url: PropTypes.string.isRequired,
    filename: PropTypes.string.isRequired,
    contentType: PropTypes.string.isRequired,
    updatedAt: PropTypes.string,
  }),
);
