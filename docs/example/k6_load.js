import http from 'k6/http';

const MAX_METRICS_PER_BATCH = 200;

export default function () {
  const namespace = 'default';
  const numberOfExperiments = 1;
  const runsPerExperiment = 2;
  const paramsPerRun = 1;
  const metricsPerRun = 2000;
  const stepsPerMetric = 4;

  for (let i = 0; i < numberOfExperiments; i++) {
    const experimentId = createExperiment(namespace);
    for (let j = 0; j < runsPerExperiment; j++) {
      createRun(namespace, experimentId, paramsPerRun, metricsPerRun, stepsPerMetric);
    }
  }
}

function createExperiment(namespace) {
  const base_url = `http://localhost:5000/ns/${namespace}/api/2.0/mlflow/`;

  const exp_response = http.post(
    base_url + 'experiments/create',
    JSON.stringify({
      "name": `experiment-${Date.now()}`,
    }),
    {
      headers: {
        'Content-Type': 'application/json'
      },
    }
  );
  return exp_response.json().experiment_id;
}

function createRun(namespace, experimentId, numParams, numMetrics, numSteps) {
  const base_url = `http://localhost:5000/ns/${namespace}/api/2.0/mlflow/`;

  const run_response = http.post(
    base_url + 'runs/create',
    JSON.stringify({
      experiment_id: experimentId,
      start_time: Date.now(),
      tags: [
        {
          key: "mlflow.user",
          value: "k6"
        }
      ]
    }),
    {
      headers: {
        'Content-Type': 'application/json'
      },
    }
  );
  const run_id = run_response.json().run.info.run_id;

  let params = [];
  for (let id = 1; id <= numParams; id++) {
    params.push({
      key: `param${id}`,
      value: `${id * Math.random()}`,
    });
  }
  http.post(
    base_url + 'runs/log-batch',
    JSON.stringify({
      run_id: run_id,
      params: params
    }),
    {
      headers: {
        'Content-Type': 'application/json'
      },
    }
  );

  let metrics = [];
  for (let step = 1; step <= numSteps; step++) {
    for (let id = 1; id <= numMetrics; id++) {
      let ctx = {};
      let rnd = Math.random();
      if (rnd < 0.3) {
        ctx = { type: 'training' };
      } else if (rnd > 0.6) {
        ctx = { type: 'testing' };
      }

      metrics.push({
        key: `metric${id}`,
        value: id * step * Math.random(),
        timestamp: Date.now(),
        step: step,
        context: ctx,
      });

      if (metrics.length >= MAX_METRICS_PER_BATCH) {
        http.post(
          base_url + 'runs/log-batch',
          JSON.stringify({
            run_id: run_id,
            metrics: metrics
          }),
          {
            headers: {
              'Content-Type': 'application/json'
            },
          }
        );
        metrics.length = 0;
      }
    }

    if (metrics.length > 0) {
      http.post(
        base_url + 'runs/log-batch',
        JSON.stringify({
          run_id: run_id,
          metrics: metrics
        }),
        {
          headers: {
            'Content-Type': 'application/json'
          },
        }
      );
    }

    logImageArtifact(namespace, run_id);

    http.post(
      base_url + 'runs/update',
      JSON.stringify({
        run_id: run_id,
        end_time: Date.now(),
        status: 'FINISHED'
      }),
      {
        headers: {
          'Content-Type': 'application/json'
        },
      }
    );
  }
}

function logImageArtifact(namespace, runId) {
  const base_url = `http://localhost:5000/ns/${namespace}/api/2.0/mlflow/`;
  
  const imageLogRequest = {
    name: 'example image name',
    iter: 1,
    step: 1,
    caption: 'example image caption',
    index: 0,
    width: 1024,
    run_id: runId,
    height: 768,
    format: 'png',
    blob_uri: 'example_image.png'
  };

  http.post(
    base_url + 'runs/log-artifact',
    JSON.stringify(imageLogRequest),
    {
      headers: {
        'Content-Type': 'application/json'
      },
    }
  );
}
