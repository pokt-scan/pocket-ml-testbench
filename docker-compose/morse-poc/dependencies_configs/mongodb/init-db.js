db = db.getSiblingDB('pocket-ml-testbench');

db.createCollection('tokenizers');
db.tokenizers.createIndex({hash: 1});

db.createCollection('tasks');
db.tasks.createIndex({
    "tasks": 1,
    "framework": 1,
    "requester_args.address": 1,
    "requester_args.service": 1,
    done: 1
});
db.createCollection('instances');
db.instances.createIndex({task_id: 1, done: 1});

db.createCollection('prompts');
db.prompts.createIndex({task_id: 1, instance_id: 1, done: 1});

db.createCollection('responses');
db.responses.createIndex({task_id: 1, instance_id: 1, prompt_id: 1, ok: 1});

db.createCollection('nodes');
db.nodes.createIndex({address: 1, service: 1}, {unique: true});

db.createCollection('results');