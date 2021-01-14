import React from 'react';
import * as PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from '../MoveDetailTable.module.scss';

import { formatWeight, formatDaysInTransit } from 'shared/formatters';

const AllowancesTable = ({ showEditBtn, info }) => {
  const titleCase = (input) => {
    if (input && input.length > 0) {
      const friendlyInput = input.toLowerCase().replace('_', ' ').split(' ');
      return friendlyInput
        .map((word) => {
          return word.replace(word[0], word[0].toUpperCase());
        })
        .join(' ');
    }
    return input;
  };
  const friendlyRankDisplay = (rank) => {
    if (rank) {
      const friendlyRank = rank.split('_');
      return `${friendlyRank[0]}-${friendlyRank[1]} ${titleCase(friendlyRank.slice(2).join(' '))}`;
    }
    return rank;
  };

  return (
    <div className={styles.MoveDetailTable}>
      <div className="stackedtable-header">
        <div>
          <h4>Allowances</h4>
        </div>
        {showEditBtn && (
          <div>
            <Link className="usa-button usa-button--secondary" data-testid="edit-allowances" to="allowances">
              Edit Allowances
            </Link>
          </div>
        )}
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Branch, rank</th>
            <td data-testid="branchRank">{`${titleCase(info.branch)}, ${friendlyRankDisplay(info.rank)}`}</td>
          </tr>
          <tr>
            <th scope="row">Weight allowance</th>
            <td data-testid="weightAllowance">{formatWeight(info.weightAllowance)}</td>
          </tr>
          <tr>
            <th scope="row">Authorized weight</th>
            <td data-testid="authorizedWeight">{formatWeight(info.authorizedWeight)}</td>
          </tr>
          <tr>
            <th scope="row">Pro-gear</th>
            <td data-testid="progear">{formatWeight(info.progear)}</td>
          </tr>
          <tr>
            <th scope="row">Spouse pro-gear</th>
            <td data-testid="spouseProgear">{formatWeight(info.spouseProgear)}</td>
          </tr>
          <tr>
            <th scope="row">Storage in transit</th>
            <td data-testid="storageInTransit">
              {info.storageInTransit ? formatDaysInTransit(info.storageInTransit) : ''}
            </td>
          </tr>
          <tr>
            <th scope="row">Dependents</th>
            <td data-testid="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</td>
          </tr>
          <tr>
            <th className="error" scope="row">
              TAC / MDC
            </th>
            <td data-testid="dependents">Missing</td>
          </tr>
          <tr>
            <th className="error" scope="row">
              SAC / SDN
            </th>
            <td data-testid="dependents">Missing</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

AllowancesTable.propTypes = {
  showEditBtn: PropTypes.bool,
  info: PropTypes.shape({
    branch: PropTypes.string,
    rank: PropTypes.string,
    weightAllowance: PropTypes.number,
    authorizedWeight: PropTypes.number,
    progear: PropTypes.number,
    spouseProgear: PropTypes.number,
    storageInTransit: PropTypes.number,
    dependents: PropTypes.bool,
  }).isRequired,
};

AllowancesTable.defaultProps = {
  showEditBtn: false,
};

export default AllowancesTable;
