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
      <div key={`${violation.id}-accordion-option`} className={styles.accordionOption}>
        <div className={styles.flex}>
          <Checkbox
            id={`${violation.id}-checkbox`}
            name={`${violation.paragraphNumber} ${violation.title}`}
            className={styles.checkbox}
            label=""
            aria-labelledby={`${violation.id}-checkbox-label`}
          />
          <div className={styles.grow} id={`${violation.id}-checkbox-label`}>
            <h5 className={styles.checkboxLabel}>{`${violation.paragraphNumber} ${violation.title}`}</h5>
            <small>{violation.requirementSummary}</small>
          </div>
          {expandedViolations.includes(violation.id) ? (
            <FontAwesomeIcon
              icon="chevron-down"
              className={styles.detailIcon}
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
              fontSize="20px"
            />
          ) : (
            <FontAwesomeIcon
              icon="chevron-up"
              className={styles.detailIcon}
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
              fontSize="20px"
            />
          )}
        </div>

        {/* Expandable Requirements Statement */}
        {expandedViolations.includes(violation.id) && (
          <p className={styles.requirementStatement}>
            <small>{violation.requirementStatement}</small>
          </p>
        )}
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
