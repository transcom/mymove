import React from 'react';
import * as PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from '../MoveDetailTable.module.scss';

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
            <th scope="row" className="text-bold">
              Branch, rank
            </th>
            <td data-testid="branchRank">{`${titleCase(info.branch)}, ${friendlyRankDisplay(info.rank)}`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Weight allowance
            </th>
            <td data-testid="weightAllowance">{`${info.weightAllowance} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Authorized weight
            </th>
            <td data-testid="authorizedWeight">{`${info.authorizedWeight} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Pro-gear
            </th>
            <td data-testid="progear">{`${info.progear} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Spouse pro-gear
            </th>
            <td data-testid="spouseProgear">{`${info.spouseProgear} lbs`}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Storage in transit
            </th>
            <td data-testid="storageInTransit">{info.storageInTransit ? `${info.storageInTransit} days` : ''}</td>
          </tr>
          <tr>
            <th scope="row" className="text-bold">
              Dependents
            </th>
            <td data-testid="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</td>
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
