-- Seed exercises data
INSERT INTO exercises (name, description, category, muscle_group, created_at, updated_at)
VALUES
    -- Chest exercises
    ('Bench Press', 'Lie on a flat bench with a barbell and press the weight up from your chest.', 'strength', 'chest', NOW(), NOW()),
    ('Push-Ups', 'A bodyweight exercise performed face-down raising and lowering your body using your arms.', 'strength', 'chest', NOW(), NOW()),
    ('Dumbbell Flyes', 'Lie on a bench with dumbbells and open your arms wide, then bring them back together.', 'strength', 'chest', NOW(), NOW()),
    ('Incline Bench Press', 'Bench press performed on an incline bench to target the upper chest.', 'strength', 'chest', NOW(), NOW()),
    ('Decline Bench Press', 'Bench press performed on a decline bench to target the lower chest.', 'strength', 'chest', NOW(), NOW()),
    
    -- Back exercises
    ('Pull-Ups', 'Hanging from a bar and pulling yourself up until your chin is over the bar.', 'strength', 'back', NOW(), NOW()),
    ('Deadlift', 'A weight training exercise where a loaded barbell is lifted from the ground to hip level.', 'strength', 'back', NOW(), NOW()),
    ('Bent Over Row', 'Bending at the hips and pulling a weight towards your lower chest/abdomen.', 'strength', 'back', NOW(), NOW()),
    ('Lat Pulldown', 'Using a cable machine to pull a bar down towards your chest.', 'strength', 'back', NOW(), NOW()),
    ('T-Bar Row', 'Using a T-bar row machine or setting up a barbell in a corner to perform rows.', 'strength', 'back', NOW(), NOW()),
    
    -- Leg exercises
    ('Squat', 'A compound exercise where you lower your body by bending your knees and hips.', 'strength', 'legs', NOW(), NOW()),
    ('Leg Press', 'Using a leg press machine to push weight away from your body with your legs.', 'strength', 'legs', NOW(), NOW()),
    ('Lunges', 'Taking a step forward and lowering your body by bending both knees.', 'strength', 'legs', NOW(), NOW()),
    ('Leg Extensions', 'Using a machine to extend your legs from a seated position.', 'strength', 'legs', NOW(), NOW()),
    ('Leg Curls', 'Using a machine to curl your legs towards your buttocks.', 'strength', 'legs', NOW(), NOW()),
    ('Calf Raises', 'Rising up on your toes to strengthen your calf muscles.', 'strength', 'legs', NOW(), NOW()),
    
    -- Shoulder exercises
    ('Overhead Press', 'Pressing a weight overhead from shoulder height.', 'strength', 'shoulders', NOW(), NOW()),
    ('Lateral Raises', 'Raising dumbbells out to the sides to target the lateral deltoids.', 'strength', 'shoulders', NOW(), NOW()),
    ('Front Raises', 'Raising dumbbells in front of you to target the anterior deltoids.', 'strength', 'shoulders', NOW(), NOW()),
    ('Reverse Flyes', 'Bending over and raising dumbbells out to the sides to target the posterior deltoids.', 'strength', 'shoulders', NOW(), NOW()),
    ('Shrugs', 'Lifting your shoulders up towards your ears while holding weights.', 'strength', 'shoulders', NOW(), NOW()),
    
    -- Arm exercises
    ('Bicep Curls', 'Curling a weight towards your shoulder to target the biceps.', 'strength', 'arms', NOW(), NOW()),
    ('Tricep Pushdowns', 'Using a cable machine to push a handle down to target the triceps.', 'strength', 'arms', NOW(), NOW()),
    ('Hammer Curls', 'Bicep curls with a neutral grip (palms facing each other).', 'strength', 'arms', NOW(), NOW()),
    ('Skull Crushers', 'Lying on a bench and lowering a weight to your forehead, then extending your arms.', 'strength', 'arms', NOW(), NOW()),
    ('Wrist Curls', 'Curling a weight using just your wrists to strengthen the forearms.', 'strength', 'arms', NOW(), NOW()),
    
    -- Core exercises
    ('Crunches', 'Lying on your back and lifting your shoulders off the ground to target the abs.', 'strength', 'core', NOW(), NOW()),
    ('Plank', 'Holding a position similar to a push-up for time to strengthen the core.', 'strength', 'core', NOW(), NOW()),
    ('Russian Twists', 'Sitting on the ground and twisting your torso from side to side.', 'strength', 'core', NOW(), NOW()),
    ('Leg Raises', 'Lying on your back and lifting your legs to target the lower abs.', 'strength', 'core', NOW(), NOW()),
    ('Mountain Climbers', 'Starting in a plank position and bringing your knees to your chest alternately.', 'strength', 'core', NOW(), NOW()),
    
    -- Cardio exercises
    ('Running', 'Moving at a pace faster than walking for an extended period.', 'cardio', 'full body', NOW(), NOW()),
    ('Cycling', 'Riding a bicycle or using a stationary bike.', 'cardio', 'legs', NOW(), NOW()),
    ('Rowing', 'Using a rowing machine to simulate the action of rowing a boat.', 'cardio', 'full body', NOW(), NOW()),
    ('Jumping Rope', 'Swinging a rope around your body and jumping over it.', 'cardio', 'full body', NOW(), NOW()),
    ('Swimming', 'Moving through water using your arms and legs.', 'cardio', 'full body', NOW(), NOW()),
    ('Elliptical Training', 'Using an elliptical machine for a low-impact cardio workout.', 'cardio', 'full body', NOW(), NOW()),
    
    -- Flexibility exercises
    ('Hamstring Stretch', 'Stretching the back of your thighs.', 'flexibility', 'legs', NOW(), NOW()),
    ('Quad Stretch', 'Stretching the front of your thighs.', 'flexibility', 'legs', NOW(), NOW()),
    ('Chest Stretch', 'Stretching your chest muscles.', 'flexibility', 'chest', NOW(), NOW()),
    ('Shoulder Stretch', 'Stretching your shoulder muscles.', 'flexibility', 'shoulders', NOW(), NOW()),
    ('Hip Flexor Stretch', 'Stretching the muscles at the front of your hips.', 'flexibility', 'legs', NOW(), NOW()),
    ('Lower Back Stretch', 'Stretching the muscles in your lower back.', 'flexibility', 'back', NOW(), NOW()); 