import PropTypes from 'prop-types';

export const FileShape = PropTypes.shape({
  filename: PropTypes.node.isRequired,
  filePath: PropTypes.string.isRequired,
  fileType: PropTypes.string.isRequired,
});

export const FilesShape = PropTypes.arrayOf(FileShape);
