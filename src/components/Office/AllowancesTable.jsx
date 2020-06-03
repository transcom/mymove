import React from 'react';
import * as PropTypes from 'prop-types';

const AllowancesTable = ({ info }) => {
  return (
    <div>
      <div className="stackedtable-header">
        <div>
          <h4>Allowances</h4>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Branch, rank</th>
            <td data-cy="branchRank">{`${info.branch}, ${info.rank}`}</td>
          </tr>
          <tr>
            <th scope="row">Weight allowance</th>
            <td data-cy="weightAllowance">{info.weightAllowance}</td>
          </tr>
          <tr>
            <th scope="row">Authorized weight</th>
            <td data-cy="authorizedWeight">{info.authorizedWeight}</td>
          </tr>
          <tr>
            <th scope="row">Pro-gear</th>
            <td data-cy="progear">{info.progear}</td>
          </tr>
          <tr>
            <th scope="row">Spouse pro-gear</th>
            <td data-cy="spouseProgear">{info.spouseProgear}</td>
          </tr>
          <tr>
            <th scope="row">Storage in transit</th>
            <td data-cy="storageInTransit">{info.storageInTransit}</td>
          </tr>
          <tr>
            <th scope="row">Dependents</th>
            <td data-cy="dependents">{info.dependents}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

AllowancesTable.propTypes = {
  info: PropTypes.shape({
    branch: PropTypes.string,
    rank: PropTypes.string,
    weightAllowance: PropTypes.string,
    authorizedWeight: PropTypes.string,
    progear: PropTypes.string,
    spouseProgear: PropTypes.string,
    storageInTransit: PropTypes.string,
    dependents: PropTypes.string,
  }).isRequired,
};

export default AllowancesTable;
