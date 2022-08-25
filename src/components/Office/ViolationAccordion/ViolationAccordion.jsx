import React from 'react';
import PropTypes from 'prop-types';
import { Grid, Accordion, Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ViolationAccordion.module.scss';

const ViolationAccordion = ({ category, violations }) => {
  const [expandedViolations, setExpandedViolations] = React.useState([]);
  const subCategories = [...new Set(violations.map((item) => item.subCategory))];

  const toggleDetailExpand = (violationId) => {
    if (expandedViolations.includes(violationId)) {
      setExpandedViolations(expandedViolations.filter((id) => id !== violationId));
    } else {
      setExpandedViolations([...expandedViolations, violationId]);
    }
  };

  const getContentForItem = (subCategory) => {
    const subCategoryViolations = violations.filter((violation) => violation.subCategory === subCategory);
    const items = subCategoryViolations.map((violation) => (
      <div style={{ borderBottom: '1px solid #F0F0F0' }} key={`${violation.id}-accordion-option`}>
        <div style={{ display: 'flex' }}>
          <Checkbox
            id={`${violation.id}-checkbox`}
            name={`${violation.paragraphNumber} ${violation.title}`}
            // label={`${violation.paragraphNumber} ${violation.title}`}
          />
          <div style={{ flexGrow: 1 }}>
            <h5 style={{ marginTop: '12px', marginBottom: 0 }}>{`${violation.paragraphNumber} ${violation.title}`}</h5>
            <small>{violation.requirementSummary}</small>
          </div>
          {expandedViolations.includes(violation.id) ? (
            <FontAwesomeIcon
              icon="chevron-down"
              style={{ color: '#565C65', fontSize: '20px', marginTop: '12px', marginRight: '12px', cursor: 'pointer' }}
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
            />
          ) : (
            <FontAwesomeIcon
              icon="chevron-up"
              style={{ color: '#565C65', fontSize: '20px', marginTop: '12px', marginRight: '12px', cursor: 'pointer' }}
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
            />
          )}
        </div>
        {expandedViolations.includes(violation.id) ? (
          <p
            style={{ marginLeft: '32px', marginTop: '8px', marginRight: '8px', marginBottom: '8px', color: '#71767A' }}
          >
            <small>{violation.requirementStatement}</small>
          </p>
        ) : null}
      </div>
    ));

    return items;
  };

  const getAccordionItems = () => {
    const items = [];
    subCategories.forEach((subCategory) => {
      items.push({
        title: subCategory,
        content: getContentForItem(subCategory),
        expanded: false,
        id: `${subCategory}-violation`,
        headingLevel: 'h4',
      });
    });

    return items;
  };

  return (
    <>
      <Grid row key={`${category}-category`}>
        <Grid col>
          <h3>{category}</h3>
        </Grid>
      </Grid>
      <div>
        <Accordion items={getAccordionItems()} multiselectable bordered className={styles.accordion} />
      </div>
    </>
  );
};

ViolationAccordion.propTypes = {
  violations: PropTypes.arrayOf(PropTypes.object),
  category: PropTypes.string,
};

ViolationAccordion.defaultProps = {
  violations: [],
  category: '',
};

export default ViolationAccordion;
