package store

// seed pre-populates the store with sample data.
func (s *Store) seed() {
	s.genres = []string{
		"Fantasy", "Action", "Romance", "Adventure", "Sci-Fi", "Mystery",
		"Horror", "Comedy", "Drama", "Slice of Life", "Martial Arts",
		"Isekai", "Wuxia", "Xianxia",
	}

	// Admin user — admin@novelhive.com / Admin123!
	admin := &User{
		ID: "admin-001", Username: "admin", Email: "admin@novelhive.com",
		PasswordHash: hashPassword("Admin123!"), Role: "admin",
	}
	s.users[admin.ID] = admin
	s.userByEmail[admin.Email] = admin.ID

	// Reader user — reader@novelhive.com / Reader123!
	reader := &User{
		ID: "reader-001", Username: "bookworm42", Email: "reader@novelhive.com",
		PasswordHash: hashPassword("Reader123!"), Role: "reader",
	}
	s.users[reader.ID] = reader
	s.userByEmail[reader.Email] = reader.ID

	// --- Novels ---
	novels := []*Novel{
		{ID: "n1", Title: "The Beginning After The End", Slug: "the-beginning-after-the-end", Author: "TurtleMe",
			Synopsis: "King Grey has unrivaled strength, wealth, and prestige in a world governed by martial ability. Yet strength is boring when it stands at the top, and lonely at that. Beneath the glamour of a powerful king lies the shell of a man, bereft of purpose and will. Reincarnated into a new world filled with magic and monsters, the king has a second chance to relive his life.",
			CoverURL: "", Status: "ongoing", TotalChapters: 3, Genres: []string{"Fantasy", "Action", "Isekai"}, ViewCount: 124500},
		{ID: "n2", Title: "Solo Leveling", Slug: "solo-leveling", Author: "Chugong",
			Synopsis: "In a world where awakened beings called Hunters fight monsters to protect humanity, Sung Jin-Woo, the weakest of all the Hunters, encounters a hidden dungeon with the hardest difficulty within the D-rank dungeons. There he finds the key to a power that could change everything.",
			CoverURL: "", Status: "completed", TotalChapters: 3, Genres: []string{"Action", "Fantasy", "Adventure"}, ViewCount: 98200},
		{ID: "n3", Title: "Omniscient Reader's Viewpoint", Slug: "omniscient-readers-viewpoint", Author: "sing N song",
			Synopsis: "Only I know the end of this world. One day our MC finds himself stuck in the world of his favorite novel. What does he do to survive? It is a unique and original story about a man who had spent most of his life reading a webnovel, and suddenly finds himself transported to that world.",
			CoverURL: "", Status: "completed", TotalChapters: 2, Genres: []string{"Fantasy", "Action", "Drama"}, ViewCount: 87600},
		{ID: "n4", Title: "Reverend Insanity", Slug: "reverend-insanity", Author: "Gu Zhen Ren",
			Synopsis: "A story of a demonic cultivator who pursues eternal life with unyielding determination. Fang Yuan is reborn 500 years into the past, armed with the knowledge of the future and a ruthless mindset. In a world of Gu Masters, he will stop at nothing to achieve his goals.",
			CoverURL: "", Status: "completed", TotalChapters: 2, Genres: []string{"Xianxia", "Fantasy", "Drama"}, ViewCount: 76300},
		{ID: "n5", Title: "Lord of the Mysteries", Slug: "lord-of-the-mysteries", Author: "Cuttlefish That Loves Diving",
			Synopsis: "In a world of mystery and steampunk, Klein Moretti finds himself reborn into a Victorian-era city. Armed with knowledge from Earth and a mysterious gray fog, he navigates a world of secret organizations, eldritch powers, and the terrifying unknown.",
			CoverURL: "", Status: "completed", TotalChapters: 2, Genres: []string{"Mystery", "Fantasy", "Adventure"}, ViewCount: 91400},
		{ID: "n6", Title: "Mushoku Tensei", Slug: "mushoku-tensei", Author: "Rifujin na Magonote",
			Synopsis: "A man reincarnated into a fantasy world as a baby with memories of his past life. Determined to live without regrets, he makes the most of his new magical abilities while navigating the complexities of family, friendship, and love in a world filled with swords and sorcery.",
			CoverURL: "", Status: "completed", TotalChapters: 2, Genres: []string{"Fantasy", "Isekai", "Adventure"}, ViewCount: 68900},
	}

	for _, n := range novels {
		s.novels[n.ID] = n
		s.novelBySlug[n.Slug] = n.ID
	}

	// --- Chapters ---
	s.chapters["n1"] = []*Chapter{
		{ID: "ch-n1-1", NovelID: "n1", Number: 1, Title: "The End of a King", WordCount: 620, Content: `The fruit of my labor was the smell of iron. The fruit of my sacrifice was the taste of blood.

I had conquered the world as King Grey, reaching the pinnacle of power. Yet, at the summit, I found nothing but emptiness—a hollow throne for a hollow king.

As my vision blurred and darkness crept in, I couldn't help but wonder: was all of this worth it? The betrayals, the wars, the countless lives lost... all for what? A title? Power?

My consciousness faded, but instead of the void I expected, I felt warmth. A strange, comforting warmth that wrapped around me like a cocoon. When I opened my eyes again, I saw nothing but white light.

And then, a voice—gentle, yet commanding—spoke words that would change everything:

"You have been given a second chance."

I blinked, confused. The light dimmed, and I found myself in a place I didn't recognize. My hands were small. My body was weak. I was... a baby?

The irony wasn't lost on me. The most powerful king in the world, reduced to a helpless infant. But as I looked up at the faces above me—faces filled with genuine love and warmth—I felt something I hadn't felt in decades.

Hope.

Perhaps this time, I could do things differently. Perhaps this time, strength wouldn't be the only thing I pursued. Perhaps this time... I could find what truly mattered.

My new mother held me close, tears streaming down her face. "Welcome to the world, Arthur," she whispered.

Arthur. That was my name now. And this was my new beginning.

As the days passed, I began to understand more about this strange new world. Magic—real, tangible magic—flowed through everything. The air itself hummed with mana, a concept I had only theorized about in my previous life.

My parents, Reynolds and Alice Leywin, were kind souls who lived in a modest home at the edge of a forest. They weren't wealthy, but they were rich in something I had never truly possessed: genuine happiness.`},
		{ID: "ch-n1-2", NovelID: "n1", Number: 2, Title: "The Beginning", WordCount: 380, Content: `Three years had passed since my rebirth into this world. In that time, I had learned to walk, talk, and—most importantly—sense mana.

The flow of mana in this world was unlike anything I had imagined. It permeated every living thing, from the smallest blade of grass to the towering ancient trees that bordered our village.

My father, Reynolds, was an adventurer by trade. He would often come home with tales of dungeons conquered and beasts slain. My mother, Alice, was a former emitter—a mage specializing in long-range magical attacks.

"Arthur, come here," my mother called from the kitchen. "Your father brought something for you."

I toddled over, maintaining the facade of a curious three-year-old, though internally I was already analyzing the magical signature of whatever lay in my father's outstretched hands.

It was a small crystal, barely the size of a marble, that pulsed with a faint blue light. A mana crystal. Even I, with all my knowledge from a previous life, was amazed at the purity of the energy contained within.

"Pretty!" I exclaimed, playing my part.

My father laughed. "Found it in the Dire Tombs. The guild said it's a D-rank mana crystal. Figured our little genius might like it."

Little did they know just how much I would indeed 'like' it. This crystal would become the foundation of my training—a tool to accelerate my understanding of mana manipulation far beyond what any child my age should be capable of.`},
		{ID: "ch-n1-3", NovelID: "n1", Number: 3, Title: "A New World", WordCount: 350, Content: `By the age of four, I had secretly mastered the basics of mana rotation—a technique that would normally take an adult mage years of dedicated practice.

The key insight from my previous life was discipline. Where young mages in this world fumbled with unfocused meditation, I applied the mental fortitude forged through decades of combat and leadership.

Every night, after my parents tucked me in, I would sit cross-legged on my bed and cycle mana through my body's core. The sensation was extraordinary—like liquid fire and cool water simultaneously flowing through invisible channels beneath my skin.

My mana core, which every person in this world possessed, was developing at an unprecedented rate. While most children didn't even become aware of their cores until age ten, mine had already reached the dark red stage.

But I knew better than to reveal this to anyone. In my previous life, prodigies were either worshipped or hunted. I suspected this world was no different.

"Arthur! Time for breakfast!" my mother's voice rang through our small home.

I hopped off the bed and padded downstairs, carefully suppressing any ambient mana leakage. To the world, I was just a bright, curious four-year-old.

To myself, I was a king in training.

The smell of freshly baked bread and forest berry jam filled the kitchen. My father sat at the table, cleaning his adventuring gear—a worn but sturdy short sword and a leather chest piece that had seen better days.

"Morning, champ," he grinned, ruffling my hair.

Today would be the day everything changed. Today, a group of bandits would attack our village, and I would have no choice but to reveal a fraction of what I was truly capable of.`},
	}

	s.chapters["n2"] = []*Chapter{
		{ID: "ch-n2-1", NovelID: "n2", Number: 1, Title: "The Weakest Hunter", WordCount: 290, Content: `Ten years ago, a portal connecting the real world to a dungeon full of monsters appeared. Ordinary people who gained the power to hunt these monsters became known as "Hunters."

However, not all Hunters are equal. While some possess incredible power capable of taking down S-rank monsters, others struggle with even the weakest E-rank creatures.

Sung Jin-Woo was the latter—an E-rank Hunter known as "humanity's weakest." Each raid was a gamble with his life, and the pay was barely enough to cover his mother's hospital bills.

"Hey, weakest! Try not to get in our way today," laughed a fellow Hunter as they prepared to enter a D-rank dungeon.

Jin-Woo gritted his teeth but said nothing. He needed the money. His younger sister was still in school, and his mother's condition wasn't improving.

The dungeon entrance shimmered like a vertical pool of blue water. One by one, the raid party stepped through.

What awaited them inside would change Jin-Woo's life forever.`},
		{ID: "ch-n2-2", NovelID: "n2", Number: 2, Title: "The Double Dungeon", WordCount: 350, Content: `Inside the D-rank dungeon, the raid proceeded normally at first. The monsters were weak—giant insects and slimes that even Jin-Woo could handle with his basic dagger.

But then they found it. Hidden behind a crumbling wall, a second entrance pulsed with ominous red light. A dungeon within a dungeon.

"We should report this to the Association," one of the senior Hunters said, his face pale.

"Are you kidding? Think of the loot!" another argued. "These hidden dungeons always have the best rewards."

Against better judgment, the group voted to enter. The moment they crossed the threshold, everything changed.

The door sealed behind them. Massive stone statues lined a cathedral-like chamber, each one depicting a warrior in a different combat pose. At the center stood a throne, and seated upon it was a figure that radiated power beyond anything they had ever encountered.

One by one, the statues began to move.

"RUN!" someone screamed, but there was nowhere to run. The sealed door would not budge.

Jin-Woo watched in horror as his companions fell one after another. The statues were merciless, their stone blades cleaving through armor and flesh alike. This wasn't a D-rank dungeon. This was death itself.

As the last few survivors huddled together, Jin-Woo noticed something the others hadn't—inscriptions on the floor. Rules. This place had rules, and following them was the only way out.

"Everyone, listen to me!" he shouted, his voice cracking with desperation. For the first time, the 'weakest' Hunter took command.`},
		{ID: "ch-n2-3", NovelID: "n2", Number: 3, Title: "Arise", WordCount: 200, Content: `Jin-Woo opened his eyes to the sterile white ceiling of a hospital room. He was alive. Somehow, impossibly, he was alive.

A translucent blue window floated before his eyes, visible only to him:

[You have been chosen as a Player of the System.]
[Daily Quest: 100 push-ups, 100 sit-ups, 100 squats, 10km run]
[Failure to complete will result in an appropriate penalty.]

He blinked. The window didn't disappear.

"What... is this?"

A new power had awakened within him—one that would transform the weakest Hunter into something the world had never seen before.`},
	}

	s.chapters["n3"] = []*Chapter{
		{ID: "ch-n3-1", NovelID: "n3", Number: 1, Title: "Prologue - Three Ways to Survive the Apocalypse", WordCount: 250, Content: `There are three ways to survive in a ruined world. I have read all of them.

For 3,149 chapters, I read a web novel called "Three Ways to Survive in a Ruined World." It took me over a decade of daily reading, but I finished every single chapter, every single word.

I was the only reader who made it to the end.

The author and I had a strange relationship. We never met, never spoke, yet through 3,149 chapters, we shared something profound. The story became my second life.

Then one morning, the subway shook. The lights flickered. Passengers screamed as cracks appeared in reality itself, and through those cracks, creatures of nightmare poured into Seoul.

The apocalypse had begun.

But unlike the seven million other people in this city, I knew exactly what was happening. Because I had read this story before.

"Scenario number one has begun," a robotic voice announced. The same words that opened Chapter 1 of the novel I had spent a decade reading.

My name is Kim Dokja, and only I know the end of this world.`},
		{ID: "ch-n3-2", NovelID: "n3", Number: 2, Title: "Scenario No.1", WordCount: 210, Content: `The subway car descended into chaos. People trampled each other trying to reach the exits, but the doors wouldn't open. Outside the windows, the sky had turned an impossible shade of crimson.

I remained seated, my worn copy of "Three Ways to Survive" clutched in my hands. My heart pounded, but my mind was clear.

According to the novel, the first scenario would test humanity's basic survival instincts. Those who panicked would die first. Those who fought recklessly would die second. Only those who understood the rules would survive.

And I knew all the rules.

"Attention passengers," the robotic voice continued. "Your survival probability will now be calculated based on your actions. Please remain calm."

A translucent status window appeared before each person. Most couldn't even read it through their tears and screaming.

I read mine carefully, then looked around the car. Somewhere among these terrified passengers was the protagonist of the story—the man who would become the strongest being in the world.

I just needed to find him before the first monster appeared. According to the novel, that would happen in exactly... three minutes.`},
	}

	// Seed remaining novels with 2 chapters each
	s.chapters["n4"] = []*Chapter{
		{ID: "ch-n4-1", NovelID: "n4", Number: 1, Title: "Rebirth", WordCount: 180, Content: `I am Fang Yuan, and I have lived for five hundred years. In my first life, I pursued the Dao of eternal life with single-minded determination. I refined Gu, cultivated my aperture, and built an empire of schemes and bloodshed. I was not a good person. I never claimed to be. The world is cruel, and I merely adapted to its cruelty. But at the end of my five hundred years, as I stood at the precipice of achieving Spring Autumn Cicada's power, I failed. My enemies surrounded me, and my body crumbled to dust. Yet in that final moment, the very Gu I had spent centuries refining activated. Time reversed. My consciousness hurtled backward through the centuries, and I awoke in the body of my fifteen-year-old self.`},
		{ID: "ch-n4-2", NovelID: "n4", Number: 2, Title: "The Path of Gu", WordCount: 160, Content: `The Qing Mao Mountain clan was as I remembered it—petty, hierarchical, and utterly predictable. The clan elders preached righteousness while secretly hoarding resources. The younger generation fought each other like dogs over scraps. Nothing had changed. I sat in my small room, feeling the familiar weight of a rank one initial stage aperture. Five hundred years of cultivation, reduced to nothing. But my mind—my mind retained everything. Every Gu recipe, every battle technique, every scheme I had ever conceived. In this world, knowledge was more dangerous than any weapon. And I possessed five centuries worth of it.`},
	}

	s.chapters["n5"] = []*Chapter{
		{ID: "ch-n5-1", NovelID: "n5", Number: 1, Title: "Crimson Moon", WordCount: 190, Content: `When Klein Moretti opened his eyes, he found himself in a dimly lit room that smelled of old wood and lamp oil. This was not his apartment. This was not the twenty-first century. A brown journal lay on the desk beside him, its pages filled with cramped handwriting in a language he shouldn't have been able to read—yet understood perfectly. The entry was dated: 7th Month, Epoch of the Grey Fog, Year 1349. "I cannot continue the ritual," the final entry read. "The voices are getting louder. If you are reading this, do not trust the Nighthawks. Do not trust the churches. Above all, do not look directly at the crimson moon." Klein's gaze drifted to the window. Outside, a blood-red moon hung low over a city of gas lamps and horse-drawn carriages, casting everything in an unsettling vermillion glow.`},
		{ID: "ch-n5-2", NovelID: "n5", Number: 2, Title: "The Grey Fog", WordCount: 170, Content: `Three days after waking in this strange body, Klein discovered the grey fog. It happened during a moment of intense concentration, when he tried to recall the details of his previous life. One second he was sitting at a desk; the next, he was floating in an infinite expanse of grey mist, standing at the head of an impossibly long bronze table. Twenty-two chairs surrounded the table, all empty except for one—the one he stood behind. Above the table, a pale golden light pulsed like a heartbeat. "What is this place?" Klein whispered. His voice echoed endlessly through the fog. This was the Sefirah Castle, though he didn't know that yet. In time, he would learn that this place existed above the spirit world, above the astral plane, above the domains of the gods themselves.`},
	}

	s.chapters["n6"] = []*Chapter{
		{ID: "ch-n6-1", NovelID: "n6", Number: 1, Title: "A Second Chance", WordCount: 180, Content: `I died at thirty-four, alone in my room, having wasted my entire life as a shut-in NEET. No friends, no accomplishments, no reason for anyone to remember me. As my consciousness faded, I felt nothing but regret. When awareness returned, I was surrounded by warmth and darkness. Muffled voices spoke in a language I didn't recognize, yet somehow understood. Then light—blinding, overwhelming light—as I was pulled into a world of vivid colors and impossible beauty. I was a newborn baby, held in the arms of a woman with pointed ears and gentle green eyes. An elf? "Rudeus," she cooed, stroking my cheek. "Welcome to our family." In that moment, cradled in my new mother's arms, I made a vow: I would not waste this life. Whatever this world offered, I would seize it. I would live without regrets.`},
		{ID: "ch-n6-2", NovelID: "n6", Number: 2, Title: "Magic Training", WordCount: 160, Content: `By age three, I could read the language of this world fluently. By four, I had secretly taught myself basic magic from a textbook I found in my father's study. The magic system here was intuitive—you gathered mana from the atmosphere, channeled it through your body, and shaped it with incantations. Simple in theory, difficult in practice. Most children didn't begin formal training until age seven. I had a three-year head start and thirty-four years of study habits from my previous life. Fire, water, earth—I practiced each element in secret, hiding my abilities from my parents. Not out of fear, but because I wanted to surprise them. The look on my father Paul's face when I silently cast a waterball spell was worth every hour of practice.`},
	}

	// Sample comments
	s.comments["ch-n1-1"] = []*Comment{
		{ID: "cm1", ChapterID: "ch-n1-1", UserID: "reader-001", Username: "bookworm42", Content: "Amazing first chapter! The transition from king to infant is handled so well.", Likes: 12},
		{ID: "cm2", ChapterID: "ch-n1-1", UserID: "admin-001", Username: "admin", Content: "One of the best isekai openings I've ever read. The introspection is what sets this apart.", Likes: 8},
	}

	// Library entries for reader
	s.library["reader-001"] = []*LibraryEntry{
		{UserID: "reader-001", NovelID: "n1", NovelTitle: "The Beginning After The End", NovelSlug: "the-beginning-after-the-end", Status: "reading", Progress: 2, Total: 3},
		{UserID: "reader-001", NovelID: "n2", NovelTitle: "Solo Leveling", NovelSlug: "solo-leveling", Status: "reading", Progress: 1, Total: 3},
	}
}
