import * as Yup from 'yup';
import PropTypes from 'prop-types';

export const uploadShape = Yup.object({
  id: Yup.string(),
  created_at: Yup.string(),
  bytes: Yup.number(),
  url: Yup.string(),
  filename: Yup.string(),
  content_type: Yup.string(),
});

export const ExistingUploadsShape = PropTypes.arrayOf(
  PropTypes.shape({
    id: PropTypes.string.isRequired,
    created_at: PropTypes.string.isRequired,
    bytes: PropTypes.number.isRequired,
    url: PropTypes.string.isRequired,
    filename: PropTypes.string.isRequired,
    content_type: PropTypes.string.isRequired,
  }),
);
