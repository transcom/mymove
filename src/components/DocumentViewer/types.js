import PropTypes from 'prop-types';

export const FileShape = PropTypes.shape({
  filename: PropTypes.node.isRequired,
  url: PropTypes.string.isRequired,
  contentType: PropTypes.string.isRequired,
});

export const FilesShape = PropTypes.arrayOf(FileShape);
