// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import 'scenes/SubmittedFeedback/IssueCards.css';

const IssueCards = ({ issues }) => {
  if (!issues) return <LoadingPlaceholder />;
  if (issues.length === 0) return <h2> There is no feedback at the moment! </h2>;
  return (
    <div className="issue-cards">
      {issues.map(issue => (
        <div key={issue.id} className="issue-card">
          {issue.description}
        </div>
      ))}
    </div>
  );
};

IssueCards.propTypes = {
  issues: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      description: PropTypes.string.isRequired,
    }),
  ),
};

export default IssueCards;
