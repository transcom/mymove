import React from 'react';
import { bool, string } from 'prop-types';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const StackedTable = ({ children, className, fullWidth }) => (
  <table className={cx(className, `${fullWidth && 'full-width-class'}`, 'default-table-classes')}>{children}</table>
);

StackedTable.propTypes = {
  fullWidth: bool,
  className: string,
  children: (props, propName, componentName) => {
    // eslint-disable-next-line security/detect-object-injection
    const prop = props[propName];
    let error;

    if (React.Children.count(prop) === 0) {
      error = new Error(`\`${componentName}\` requires Children.`);
    }
    React.Children.forEach(prop, (el) => {
      if (error) return;
      if (el.type.name !== 'StackedTableRow') {
        error = new Error(`\`${componentName}\` children must be \`StackedTableRow\`.`);
      }
    });

    return error;
  },
};

const StackedTableRow = ({ children, className }) => (
  <tr className={cx(className, 'default-table-row-classes')}>{children}</tr>
);

StackedTableRow.propTypes = {
  className: string,
  children: (props, propName, componentName) => {
    // eslint-disable-next-line security/detect-object-injection
    const prop = props[propName];
    let error;

    if (React.Children.count(prop) === 0) {
      error = new Error(`\`${componentName}\` requires Children.`);
    }

    React.Children.forEach(prop, (el) => {
      if (error) return;
      if (el.type.name !== 'StackedTableHeader') {
        error = new Error(`\`${componentName}\` children must be \`StackedTableHeader\`.`);
      } else if (el.type.name !== 'StackedTableData') {
        error = new Error(`\`${componentName}\` children must be \`StackedTableData\`.`);
      }
    });

    return error;
  },
};

const StackedTableHeader = ({ children, className }) => (
  <th className={cx('default-table-header-class-names', className)}>{children}</th>
);

StackedTableHeader.propTypes = {
  className: string,
  children: (props, propName, componentName) =>
    // eslint-disable-next-line security/detect-object-injection
    React.Children.count(props[propName]) === 0 ? new Error(`\`${componentName}\` requires Children.`) : null,
};

const StackedTableData = ({ children, className }) => (
  <td className={cx('default-table-data-class-names', className)}>{children}</td>
);

StackedTableData.propTypes = {
  className: string,
  children: (props, propName, componentName) =>
    // eslint-disable-next-line security/detect-object-injection
    React.Children.count(props[propName]) === 0 ? new Error(`\`${componentName}\` requires Children.`) : null,
};

export { StackedTable, StackedTableRow, StackedTableHeader, StackedTableData };
