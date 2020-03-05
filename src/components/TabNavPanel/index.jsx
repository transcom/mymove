import React from 'react';
import PropTypes from 'prop-types';
import { TabPanel } from 'react-tabs';

// eslint-disable-next-line react/jsx-props-no-spreading
const TabNavPanel = ({ children, ...otherProps }) => <TabPanel {...otherProps}>{children}</TabPanel>;

TabNavPanel.propTypes = {
  children: PropTypes.node.isRequired,
};

TabNavPanel.tabsRole = 'TabPanel';

export default TabNavPanel;
