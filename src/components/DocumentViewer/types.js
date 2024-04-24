import PropTypes from 'prop-types';

export const FileShape = PropTypes.shape({
  filename: PropTypes.node.isRequired,
  url: PropTypes.string.isRequired,
  contentType: PropTypes.string.isRequired,
  isWeightTicket: PropTypes.bool, // Optional field to determine if the file is a weight ticket or not. Used within DocViewerMenu
});

export const FilesShape = PropTypes.arrayOf(FileShape);
